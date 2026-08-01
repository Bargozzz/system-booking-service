package handlers

import (
	"log"
	"time"

	"minibooking/db"
)

// StartLockExpirySweeper runs in the background and periodically releases
// any seat locks that have passed their expiry time, so seats become
// available again even if no one queries them directly.
func StartLockExpirySweeper(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			res, err := db.DB.Exec(
				`UPDATE seats SET status = 'available', locked_by = NULL, lock_expires_at = NULL
				 WHERE status = 'locked' AND lock_expires_at < ?`,
				time.Now(),
			)
			if err != nil {
				log.Printf("lock sweeper error: %v", err)
				continue
			}
			if n, _ := res.RowsAffected(); n > 0 {
				log.Printf("lock sweeper: released %d expired seat lock(s)", n)
			}
		}
	}()
}
