package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init opens the sqlite database at path, enables WAL + busy_timeout for
// better concurrent-write behavior, and runs migrations + seed data.
func Init(path string) error {
	conn, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", path))
	if err != nil {
		return err
	}
	// SQLite only supports a single writer at a time; limiting the pool
	// avoids "database is locked" errors under concurrent access and lets
	// the busy_timeout pragma queue writers instead of failing immediately.
	conn.SetMaxOpenConns(1)

	DB = conn

	if err := migrate(); err != nil {
		return err
	}
	if err := seed(); err != nil {
		return err
	}
	return nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS operators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS routes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operator_id INTEGER NOT NULL REFERENCES operators(id),
		origin TEXT NOT NULL,
		destination TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		route_id INTEGER NOT NULL REFERENCES routes(id),
		departure_time DATETIME NOT NULL,
		arrival_time DATETIME NOT NULL,
		price REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS seats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER NOT NULL REFERENCES schedules(id),
		seat_number TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'available',
		locked_by INTEGER,
		lock_expires_at DATETIME,
		UNIQUE(schedule_id, seat_number)
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id),
		schedule_id INTEGER NOT NULL REFERENCES schedules(id),
		seat_id INTEGER NOT NULL REFERENCES seats(id),
		created_at DATETIME NOT NULL,
		status TEXT NOT NULL DEFAULT 'confirmed'
	);

	CREATE INDEX IF NOT EXISTS idx_schedules_route ON schedules(route_id);
	CREATE INDEX IF NOT EXISTS idx_seats_schedule ON seats(schedule_id);
	`
	_, err := DB.Exec(schema)
	return err
}

func seed() error {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	log.Println("Seeding database with dummy data...")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Dummy users (plaintext passwords for this simplified demo auth)
	users := []struct{ username, password string }{
		{"alice", "password123"},
		{"bob", "password123"},
	}
	for _, u := range users {
		if _, err := tx.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, u.username, u.password); err != nil {
			return err
		}
	}

	// Operators
	res, err := tx.Exec(`INSERT INTO operators (name) VALUES (?)`, "Tiketux Air")
	if err != nil {
		return err
	}
	operatorID, _ := res.LastInsertId()

	// Routes
	res, err = tx.Exec(`INSERT INTO routes (operator_id, origin, destination) VALUES (?, ?, ?)`,
		operatorID, "Jakarta", "Bandung")
	if err != nil {
		return err
	}
	routeID, _ := res.LastInsertId()

	res2, err := tx.Exec(`INSERT INTO routes (operator_id, origin, destination) VALUES (?, ?, ?)`,
		operatorID, "Bandung", "Surabaya")
	if err != nil {
		return err
	}
	routeID2, _ := res2.LastInsertId()

	now := time.Now()
	// Schedules for route 1
	sched1, err := tx.Exec(`INSERT INTO schedules (route_id, departure_time, arrival_time, price) VALUES (?, ?, ?, ?)`,
		routeID, now.Add(24*time.Hour), now.Add(26*time.Hour), 250000.0)
	if err != nil {
		return err
	}
	sched1ID, _ := sched1.LastInsertId()

	sched2, err := tx.Exec(`INSERT INTO schedules (route_id, departure_time, arrival_time, price) VALUES (?, ?, ?, ?)`,
		routeID, now.Add(48*time.Hour), now.Add(50*time.Hour), 275000.0)
	if err != nil {
		return err
	}
	sched2ID, _ := sched2.LastInsertId()

	// Schedule for route 2
	sched3, err := tx.Exec(`INSERT INTO schedules (route_id, departure_time, arrival_time, price) VALUES (?, ?, ?, ?)`,
		routeID2, now.Add(30*time.Hour), now.Add(38*time.Hour), 400000.0)
	if err != nil {
		return err
	}
	sched3ID, _ := sched3.LastInsertId()

	// Seats: 12 seats per schedule (A1-A6, B1-B6)
	for _, schedID := range []int64{sched1ID, sched2ID, sched3ID} {
		for _, row := range []string{"A", "B"} {
			for n := 1; n <= 6; n++ {
				seatNum := fmt.Sprintf("%s%d", row, n)
				if _, err := tx.Exec(`INSERT INTO seats (schedule_id, seat_number, status) VALUES (?, ?, 'available')`,
					schedID, seatNum); err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}
