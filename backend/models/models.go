package models

import "time"

// Operator represents a transportation/event operator (e.g. an airline, bus company, cinema)
type Operator struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Route represents an origin -> destination pair served by an operator
type Route struct {
	ID          int64  `json:"id"`
	OperatorID  int64  `json:"operator_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

// Schedule represents a specific departure of a route on a given date/time
type Schedule struct {
	ID            int64     `json:"id"`
	RouteID       int64     `json:"route_id"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Price         float64   `json:"price"`

	// Populated for search results (joined data)
	OperatorName string `json:"operator_name,omitempty"`
	Origin       string `json:"origin,omitempty"`
	Destination  string `json:"destination,omitempty"`
}

// SeatStatus enumerates the lifecycle of a seat for a given schedule
type SeatStatus string

const (
	SeatAvailable SeatStatus = "available"
	SeatLocked    SeatStatus = "locked"
	SeatBooked    SeatStatus = "booked"
)

// Seat represents an individual seat belonging to a schedule
type Seat struct {
	ID           int64      `json:"id"`
	ScheduleID   int64      `json:"schedule_id"`
	SeatNumber   string     `json:"seat_number"`
	Status       SeatStatus `json:"status"`
	LockedBy     *int64     `json:"locked_by,omitempty"`
	LockExpiresAt *time.Time `json:"lock_expires_at,omitempty"`
}

// Booking represents a confirmed reservation of a seat by a user
type Booking struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ScheduleID int64     `json:"schedule_id"`
	SeatID     int64     `json:"seat_id"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"` // confirmed, cancelled
}

// User represents a pre-seeded dummy user for authentication
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // stored as plaintext for this demo/simplified auth, never exposed
}
