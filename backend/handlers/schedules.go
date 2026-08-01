package handlers

import (
	"net/http"
	"time"

	"minibooking/db"
	"minibooking/models"
)

// SearchSchedules handles GET /api/schedules?origin=X&destination=Y&date=YYYY-MM-DD
// origin, destination, and date are all optional filters; omitting all
// returns every upcoming schedule.
func SearchSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	dateStr := r.URL.Query().Get("date")

	query := `
		SELECT s.id, s.route_id, s.departure_time, s.arrival_time, s.price,
		       o.name, rt.origin, rt.destination
		FROM schedules s
		JOIN routes rt ON s.route_id = rt.id
		JOIN operators o ON rt.operator_id = o.id
		WHERE 1=1
	`
	args := []interface{}{}
	if origin != "" {
		query += " AND rt.origin LIKE ?"
		args = append(args, "%"+origin+"%")
	}
	if destination != "" {
		query += " AND rt.destination LIKE ?"
		args = append(args, "%"+destination+"%")
	}
	if dateStr != "" {
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			writeJSONError(w, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
			return
		}
		query += " AND date(s.departure_time) = date(?)"
		args = append(args, dateStr)
	}
	query += " ORDER BY s.departure_time ASC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to query schedules")
		return
	}
	defer rows.Close()

	results := []models.Schedule{}
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.RouteID, &s.DepartureTime, &s.ArrivalTime, &s.Price,
			&s.OperatorName, &s.Origin, &s.Destination); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to read schedule row")
			return
		}
		results = append(results, s)
	}

	writeJSON(w, http.StatusOK, results)
}

// ViewSeats handles GET /api/schedules/{id}/seats and returns seat
// availability for a schedule, auto-releasing any expired locks first.
func ViewSeats(scheduleID int64, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Lazily release any locks that expired since the last sweep, so
	// clients always see accurate availability even between sweep ticks.
	if _, err := db.DB.Exec(
		`UPDATE seats SET status = 'available', locked_by = NULL, lock_expires_at = NULL
		 WHERE schedule_id = ? AND status = 'locked' AND lock_expires_at < ?`,
		scheduleID, time.Now(),
	); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to refresh seat status")
		return
	}

	rows, err := db.DB.Query(
		`SELECT id, schedule_id, seat_number, status, locked_by, lock_expires_at
		 FROM seats WHERE schedule_id = ? ORDER BY seat_number ASC`, scheduleID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to query seats")
		return
	}
	defer rows.Close()

	seats := []models.Seat{}
	for rows.Next() {
		var s models.Seat
		if err := rows.Scan(&s.ID, &s.ScheduleID, &s.SeatNumber, &s.Status, &s.LockedBy, &s.LockExpiresAt); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to read seat row")
			return
		}
		seats = append(seats, s)
	}
	if len(seats) == 0 {
		writeJSONError(w, http.StatusNotFound, "schedule not found or has no seats")
		return
	}

	writeJSON(w, http.StatusOK, seats)
}
