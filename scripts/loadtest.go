// Command loadtest fires N concurrent requests at the seat-lock endpoint
// for a single seat, to demonstrate that the backend's concurrency
// handling allows exactly one request to succeed no matter how many
// requests race for the same seat at once.
//
// Usage (from the mini-booking-services root, backend must already be
// running on :8080):
//
//	go run scripts/loadtest.go
//	go run scripts/loadtest.go -requests=100 -seat=5
//
// It logs in as the seeded "alice" user, discovers the first available
// seat on the first schedule it finds (or uses -seat if given), then
// fires all requests at once and prints a breakdown of status codes.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type schedule struct {
	ID int64 `json:"id"`
}

type seat struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func main() {
	baseURL := flag.String("base-url", "http://localhost:8080", "backend base URL")
	requests := flag.Int("requests", 50, "number of concurrent lock requests to fire")
	seatIDFlag := flag.Int64("seat", 0, "seat ID to target (0 = auto-pick first available seat)")
	flag.Parse()

	fmt.Println("== Mini Booking Service — concurrency load test ==")
	fmt.Printf("Target: %s | concurrent requests: %d\n\n", *baseURL, *requests)

	token := login(*baseURL, "alice", "password123")
	fmt.Println("✓ Logged in as alice")

	seatID := *seatIDFlag
	if seatID == 0 {
		seatID = findAvailableSeat(*baseURL)
	}
	fmt.Printf("✓ Targeting seat ID %d\n\n", seatID)

	fmt.Printf("Firing %d concurrent POST /api/seats/%d/lock requests...\n", *requests, seatID)

	var wg sync.WaitGroup
	results := make(chan int, *requests)
	start := time.Now()

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- lockSeat(*baseURL, seatID, token)
		}()
	}

	wg.Wait()
	close(results)
	elapsed := time.Since(start)

	tally := map[int]int{}
	for code := range results {
		tally[code]++
	}

	fmt.Printf("\nDone in %s.\n\n", elapsed)
	fmt.Println("Result breakdown:")
	for code, count := range tally {
		label := "other"
		switch code {
		case http.StatusOK:
			label = "200 OK (lock acquired)"
		case http.StatusConflict:
			label = "409 Conflict (seat already locked/booked)"
		case http.StatusUnauthorized:
			label = "401 Unauthorized"
		}
		fmt.Printf("  %-45s x%d\n", label, count)
	}

	successes := tally[http.StatusOK]
	fmt.Println()
	if successes == 1 {
		fmt.Println("✅ PASS: exactly 1 request acquired the lock, as expected.")
	} else {
		fmt.Printf("❌ FAIL: expected exactly 1 successful lock, got %d. Concurrency handling may be broken!\n", successes)
	}
}

func login(baseURL, username, password string) string {
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	resp, err := http.Post(baseURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(fmt.Sprintf("login request failed: %v", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("login failed with status %d: %s", resp.StatusCode, b))
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	token, _ := result["access_token"].(string)
	if token == "" {
		panic("login response missing access_token")
	}
	return token
}

func findAvailableSeat(baseURL string) int64 {
	resp, err := http.Get(baseURL + "/api/schedules")
	if err != nil {
		panic(fmt.Sprintf("failed to fetch schedules: %v", err))
	}
	defer resp.Body.Close()
	var schedules []schedule
	json.NewDecoder(resp.Body).Decode(&schedules)
	if len(schedules) == 0 {
		panic("no schedules found — has the backend been seeded?")
	}

	seatsResp, err := http.Get(fmt.Sprintf("%s/api/schedules/%d/seats", baseURL, schedules[0].ID))
	if err != nil {
		panic(fmt.Sprintf("failed to fetch seats: %v", err))
	}
	defer seatsResp.Body.Close()
	var seats []seat
	json.NewDecoder(seatsResp.Body).Decode(&seats)

	for _, s := range seats {
		if s.Status == "available" {
			return s.ID
		}
	}
	panic("no available seats found on the first schedule — try again after the sweeper releases expired locks, or specify -seat explicitly")
}

func lockSeat(baseURL string, seatID int64, token string) int {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/seats/%d/lock", baseURL, seatID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}
