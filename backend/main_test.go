package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"minibooking/db"
)

// setupTestServer initializes a fresh temp sqlite DB and starts an
// httptest.Server backed by the real router, so tests exercise the full
// HTTP + auth + concurrency stack exactly as production does.
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "minibooking-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp db file: %v", err)
	}
	tmpFile.Close()
	dbPath := tmpFile.Name()

	if err := db.Init(dbPath); err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	srv := httptest.NewServer(NewRouter())
	t.Cleanup(func() {
		srv.Close()
		db.DB.Close()
		os.Remove(dbPath)
		os.Remove(dbPath + "-wal")
		os.Remove(dbPath + "-shm")
	})
	return srv
}

func login(t *testing.T, baseURL, username, password string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	resp, err := http.Post(baseURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected login to succeed, got status %d", resp.StatusCode)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	token, _ := result["access_token"].(string)
	if token == "" {
		t.Fatalf("login response missing access_token")
	}
	return token
}

func doPost(baseURL, path, token string) *http.Response {
	req, _ := http.NewRequest(http.MethodPost, baseURL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	return resp
}

// TestSuccessfulBooking covers the full happy path: search -> view seats ->
// lock -> confirm.
func TestSuccessfulBooking(t *testing.T) {
	srv := setupTestServer(t)
	token := login(t, srv.URL, "alice", "password123")

	resp, err := http.Get(srv.URL + "/api/schedules?origin=Jakarta")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected schedule search to succeed, err=%v status=%v", err, resp)
	}
	var schedules []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&schedules)
	resp.Body.Close()
	if len(schedules) == 0 {
		t.Fatalf("expected at least one schedule from seeded data")
	}
	scheduleID := int64(schedules[0]["id"].(float64))

	seatsResp, err := http.Get(srv.URL + "/api/schedules/" + itoa(scheduleID) + "/seats")
	if err != nil || seatsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected seat view to succeed, err=%v", err)
	}
	var seats []map[string]interface{}
	json.NewDecoder(seatsResp.Body).Decode(&seats)
	seatsResp.Body.Close()
	if len(seats) == 0 {
		t.Fatalf("expected seats to exist for schedule")
	}
	seatID := int64(seats[0]["id"].(float64))

	lockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", token)
	if lockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected lock to succeed, got %d", lockResp.StatusCode)
	}
	lockResp.Body.Close()

	confirmResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/confirm", token)
	if confirmResp.StatusCode != http.StatusOK {
		t.Fatalf("expected confirm to succeed, got %d", confirmResp.StatusCode)
	}
	var confirmResult map[string]interface{}
	json.NewDecoder(confirmResp.Body).Decode(&confirmResult)
	confirmResp.Body.Close()
	if confirmResult["status"] != "confirmed" {
		t.Fatalf("expected booking status 'confirmed', got %v", confirmResult["status"])
	}
}

// TestLockConflict verifies a second user cannot lock or confirm a seat
// that is already validly locked by another user (race-condition safety).
func TestLockConflict(t *testing.T) {
	srv := setupTestServer(t)
	aliceToken := login(t, srv.URL, "alice", "password123")
	bobToken := login(t, srv.URL, "bob", "password123")

	const seatID = 1

	lockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", aliceToken)
	if lockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected alice's lock to succeed, got %d", lockResp.StatusCode)
	}
	lockResp.Body.Close()

	conflictResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", bobToken)
	if conflictResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected bob's lock to conflict with 409, got %d", conflictResp.StatusCode)
	}
	conflictResp.Body.Close()

	confirmAsBobResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/confirm", bobToken)
	if confirmAsBobResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected bob's confirm to fail with 409 (not his lock), got %d", confirmAsBobResp.StatusCode)
	}
	confirmAsBobResp.Body.Close()
}

// TestConcurrentLockingIsRaceFree fires many simultaneous lock requests at
// the same seat and asserts exactly one succeeds.
func TestConcurrentLockingIsRaceFree(t *testing.T) {
	srv := setupTestServer(t)
	token := login(t, srv.URL, "alice", "password123")

	const seatID = 2
	const attempts = 15
	results := make(chan int, attempts)

	for i := 0; i < attempts; i++ {
		go func() {
			resp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", token)
			results <- resp.StatusCode
			resp.Body.Close()
		}()
	}

	successCount := 0
	for i := 0; i < attempts; i++ {
		if <-results == http.StatusOK {
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful lock out of %d concurrent attempts, got %d", attempts, successCount)
	}
}

// TestMissingOrInvalidToken verifies protected endpoints reject requests
// with no token and with a garbage token.
func TestMissingOrInvalidToken(t *testing.T) {
	srv := setupTestServer(t)

	noTokenResp := doPost(srv.URL, "/api/seats/1/lock", "")
	if noTokenResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 with no token, got %d", noTokenResp.StatusCode)
	}
	noTokenResp.Body.Close()

	invalidTokenResp := doPost(srv.URL, "/api/seats/1/lock", "not.a.valid.jwt")
	if invalidTokenResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 with invalid token, got %d", invalidTokenResp.StatusCode)
	}
	invalidTokenResp.Body.Close()
}

// TestLockExpiry verifies that once a lock's expiry time has passed, the
// seat becomes lockable by another user again (simulating the sweeper /
// lazy-release behavior without waiting the full 5 minutes).
func TestLockExpiry(t *testing.T) {
	srv := setupTestServer(t)
	aliceToken := login(t, srv.URL, "alice", "password123")
	bobToken := login(t, srv.URL, "bob", "password123")

	const seatID = 3

	lockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", aliceToken)
	if lockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected alice's lock to succeed, got %d", lockResp.StatusCode)
	}
	lockResp.Body.Close()

	// Simulate time passing by directly rewinding the lock's expiry into
	// the past, rather than sleeping for the real 5-minute lock duration.
	_, err := db.DB.Exec(`UPDATE seats SET lock_expires_at = ? WHERE id = ?`,
		time.Now().Add(-1*time.Minute), seatID)
	if err != nil {
		t.Fatalf("failed to simulate lock expiry: %v", err)
	}

	// Bob should now be able to lock the seat since alice's lock expired.
	bobLockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", bobToken)
	if bobLockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected bob to acquire lock after expiry, got %d", bobLockResp.StatusCode)
	}
	bobLockResp.Body.Close()

	// Alice trying to confirm her now-expired lock should fail.
	aliceConfirmResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/confirm", aliceToken)
	if aliceConfirmResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected alice's confirm on expired lock to fail with 409, got %d", aliceConfirmResp.StatusCode)
	}
	aliceConfirmResp.Body.Close()
}

// TestInvalidLoginCredentials verifies bad credentials are rejected.
func TestInvalidLoginCredentials(t *testing.T) {
	srv := setupTestServer(t)
	body, _ := json.Marshal(map[string]string{"username": "alice", "password": "wrongpassword"})
	resp, err := http.Post(srv.URL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong password, got %d", resp.StatusCode)
	}
}

// TestCancelLock verifies a user can voluntarily release their own pending
// hold (e.g. they clicked the wrong seat), and that the seat becomes
// available to other users immediately afterwards.
func TestCancelLock(t *testing.T) {
	srv := setupTestServer(t)
	aliceToken := login(t, srv.URL, "alice", "password123")
	bobToken := login(t, srv.URL, "bob", "password123")

	const seatID = 4

	lockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", aliceToken)
	if lockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected alice's lock to succeed, got %d", lockResp.StatusCode)
	}
	lockResp.Body.Close()

	cancelResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/unlock", aliceToken)
	if cancelResp.StatusCode != http.StatusOK {
		t.Fatalf("expected alice's cancel to succeed, got %d", cancelResp.StatusCode)
	}
	cancelResp.Body.Close()

	// Bob should now be able to lock the same seat immediately.
	bobLockResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/lock", bobToken)
	if bobLockResp.StatusCode != http.StatusOK {
		t.Fatalf("expected bob to lock the seat right after alice cancelled, got %d", bobLockResp.StatusCode)
	}
	bobLockResp.Body.Close()

	// Alice can no longer confirm a lock she cancelled.
	aliceConfirmResp := doPost(srv.URL, "/api/seats/"+itoa(seatID)+"/confirm", aliceToken)
	if aliceConfirmResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected alice's confirm after cancelling to fail with 409, got %d", aliceConfirmResp.StatusCode)
	}
	aliceConfirmResp.Body.Close()
}

func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}
