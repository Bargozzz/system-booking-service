package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"minibooking/db"
	"minibooking/handlers"
	"minibooking/middleware"
	"minibooking/ratelimit"
)

// corsMiddleware allows the Vue.js dev server (or any origin, for this demo)
// to call the API from the browser.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// pathID extracts the trailing numeric ID from a URL path like
// /api/seats/42/lock -> parts[2] == "42".
func pathID(path string, position int) (int64, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if position >= len(parts) {
		return 0, false
	}
	id, err := strconv.ParseInt(parts[position], 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

// NewRouter builds the full HTTP route table. Extracted from main() so
// integration tests can spin up a real httptest.Server against it.
func NewRouter() http.Handler {
	searchLimiter := ratelimit.New(2, 10) // 2 req/sec sustained, burst of 10

	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/api/login", corsMiddleware(handlers.Login))
	mux.HandleFunc("/api/refresh", corsMiddleware(handlers.Refresh))
	mux.HandleFunc("/api/schedules", corsMiddleware(searchLimiter.Middleware(handlers.SearchSchedules)))

	// /api/schedules/{id}/seats  (public - view availability)
	mux.HandleFunc("/api/schedules/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(r.URL.Path, 2)
		if !ok || !strings.HasSuffix(r.URL.Path, "/seats") {
			http.NotFound(w, r)
			return
		}
		handlers.ViewSeats(id, w, r)
	}))

	// /api/seats/{id}/lock  (protected)
	// /api/seats/{id}/confirm  (protected)
	// /api/seats/{id}/unlock  (protected)
	mux.HandleFunc("/api/seats/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathID(r.URL.Path, 2)
		if !ok {
			http.NotFound(w, r)
			return
		}
		switch {
		case strings.HasSuffix(r.URL.Path, "/lock"):
			middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
				handlers.LockSeat(id, w, r)
			}).ServeHTTP(w, r)
		case strings.HasSuffix(r.URL.Path, "/confirm"):
			middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
				handlers.ConfirmBooking(id, w, r)
			}).ServeHTTP(w, r)
		case strings.HasSuffix(r.URL.Path, "/unlock"):
			middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
				handlers.UnlockSeat(id, w, r)
			}).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	return mux
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./minibooking.db"
	}
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	log.Println("Database initialized and seeded.")

	handlers.StartLockExpirySweeper(30 * time.Second)
	log.Println("Lock expiry sweeper started (interval: 30s).")

	router := NewRouter()

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("Mini Booking Service listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
