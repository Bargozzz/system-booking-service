package handlers

import (
	"net/http"
	"time"

	"minibooking/db"
	"minibooking/middleware"
)

// LockDuration is how long a seat lock is held before it auto-expires.
const LockDuration = 5 * time.Minute

// LockSeat handles POST /api/seats/{id}/lock (protected).
// It atomically transitions a seat from available (or expired-locked) to
// locked-by-this-user using a single conditional UPDATE, which prevents two
// concurrent requests from both succeeding for the same seat: whichever
// UPDATE commits first changes the row's state, and the second one's WHERE
// clause no longer matches, so its RowsAffected is 0.
func LockSeat(seatID int64, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	now := time.Now()
	expiresAt := now.Add(LockDuration)

	res, err := db.DB.Exec(
		`UPDATE seats
		 SET status = 'locked', locked_by = ?, lock_expires_at = ?
		 WHERE id = ?
		   AND (status = 'available' OR (status = 'locked' AND lock_expires_at < ?))`,
		userID, expiresAt, seatID, now,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to lock seat")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Either the seat doesn't exist, is already booked, or is validly
		// locked by someone else right now.
		var status string
		checkErr := db.DB.QueryRow(`SELECT status FROM seats WHERE id = ?`, seatID).Scan(&status)
		if checkErr != nil {
			writeJSONError(w, http.StatusNotFound, "seat not found")
			return
		}
		if status == "booked" {
			writeJSONError(w, http.StatusConflict, "seat is already booked")
			return
		}
		writeJSONError(w, http.StatusConflict, "seat is currently locked by another user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"seat_id":         seatID,
		"status":          "locked",
		"lock_expires_at": expiresAt,
	})
}

// UnlockSeat handles POST /api/seats/{id}/unlock (protected).
// Lets a user voluntarily release their own pending hold before it expires
// or before confirming — e.g. because they clicked the wrong seat, or want
// to pick a different one. It only releases the seat if the caller is the
// one currently holding the lock; it never touches a seat locked by someone
// else or one that's already booked.
func UnlockSeat(seatID int64, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	res, err := db.DB.Exec(
		`UPDATE seats SET status = 'available', locked_by = NULL, lock_expires_at = NULL
		 WHERE id = ? AND status = 'locked' AND locked_by = ?`,
		seatID, userID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to release seat")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Nothing to release: either it wasn't locked by this user, it
		// already expired (sweeper/lazy-release beat us to it), or it
		// doesn't exist. Any of those cases means the seat is no longer
		// "held by you", which is the state the caller wanted anyway — so
		// this is treated as a harmless no-op rather than an error.
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"seat_id": seatID,
			"status":  "released",
			"note":    "seat was not held by you (already expired or released)",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"seat_id": seatID,
		"status":  "released",
	})
}

// ConfirmBooking handles POST /api/seats/{id}/confirm (protected).
// It atomically transitions a seat from locked-by-this-user (and not yet
// expired) to booked, and inserts the booking record in the same
// transaction so the two states never drift apart.
func ConfirmBooking(seatID int64, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "missing user context")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to start transaction")
		return
	}
	defer tx.Rollback()

	now := time.Now()
	var scheduleID int64
	var status string
	var lockedBy *int64
	var lockExpiresAt *time.Time
	err = tx.QueryRow(
		`SELECT schedule_id, status, locked_by, lock_expires_at FROM seats WHERE id = ?`, seatID,
	).Scan(&scheduleID, &status, &lockedBy, &lockExpiresAt)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "seat not found")
		return
	}

	if status != "locked" || lockedBy == nil || *lockedBy != userID {
		writeJSONError(w, http.StatusConflict, "seat is not locked by you; please lock it again")
		return
	}
	if lockExpiresAt == nil || lockExpiresAt.Before(now) {
		// Expired: release it so it becomes available again, then reject.
		tx.Exec(`UPDATE seats SET status='available', locked_by=NULL, lock_expires_at=NULL WHERE id=?`, seatID)
		tx.Commit()
		writeJSONError(w, http.StatusConflict, "seat lock has expired; please lock the seat again")
		return
	}

	// Note: locked_by is intentionally NOT cleared here. Once a seat's
	// status becomes 'booked', locked_by is repurposed from "who currently
	// holds the temporary lock" to "who owns the confirmed booking" — this
	// lets clients distinguish "booked by me" from "booked by someone else"
	// without a second query.
	res, err := tx.Exec(
		`UPDATE seats SET status = 'booked', lock_expires_at = NULL
		 WHERE id = ? AND status = 'locked' AND locked_by = ?`,
		seatID, userID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to confirm booking")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		writeJSONError(w, http.StatusConflict, "seat state changed, please try locking again")
		return
	}

	bookingRes, err := tx.Exec(
		`INSERT INTO bookings (user_id, schedule_id, seat_id, created_at, status) VALUES (?, ?, ?, ?, 'confirmed')`,
		userID, scheduleID, seatID, now,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create booking record")
		return
	}
	bookingID, _ := bookingRes.LastInsertId()

	if err := tx.Commit(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to commit booking")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"booking_id":  bookingID,
		"seat_id":     seatID,
		"schedule_id": scheduleID,
		"status":      "confirmed",
	})
}
