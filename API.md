# API Documentation — Mini Booking Service

Base URL (local): `http://localhost:8080`

All request/response bodies are JSON. All examples use `curl`; import them
into Postman if you prefer a GUI (each block below is a valid standalone
request).

---

## 1. Authentication

### 1.1 Login

Get an access token + refresh token using one of the seeded dummy accounts.

**Seeded users** (no registration needed): `alice` / `password123`, `bob` / `password123`.

```
POST /api/login
Content-Type: application/json
```

Request body:
```json
{
  "username": "alice",
  "password": "password123"
}
```

Success response — `200 OK`:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "user_id": 1,
  "username": "alice"
}
```

`user_id` is returned so the frontend can compare it against a seat's
`locked_by` field to tell "booked by me" apart from "booked by someone
else" — see section 2.2 below.

Failure response — `401 Unauthorized`:
```json
{ "error": "invalid username or password" }
```

`curl` example:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

### 1.2 Using the access token

Pass it as a Bearer token in the `Authorization` header on any protected
endpoint (seat lock, booking confirmation):

```
Authorization: Bearer <access_token>
```

The access token expires after **15 minutes**. After it expires, protected
endpoints respond `401 Unauthorized` with `{"error":"invalid or expired token"}`
— at that point, use the refresh endpoint below instead of logging in again.

### 1.3 Refresh an expired access token

```
POST /api/refresh
Content-Type: application/json
```

Request body:
```json
{ "refresh_token": "<refresh_token from login>" }
```

Success response — `200 OK`:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

The refresh token itself is valid for **24 hours**. Once it expires, the
user must log in again.

`curl` example:
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<paste refresh_token here>"}'
```

---

## 2. Public endpoints (no token required)

### 2.1 Search schedules

```
GET /api/schedules?origin=Jakarta&destination=Bandung&date=2026-08-01
```

All three query params are optional filters; omit any/all to broaden the
search. Rate-limited to 2 requests/second per IP (burst of 10) — exceeding
this returns `429 Too Many Requests`.

Response — `200 OK`:
```json
[
  {
    "id": 1,
    "route_id": 1,
    "departure_time": "2026-08-01T09:00:00Z",
    "arrival_time": "2026-08-01T11:00:00Z",
    "price": 250000,
    "operator_name": "Tiketux Air",
    "origin": "Jakarta",
    "destination": "Bandung"
  }
]
```

`curl` example:
```bash
curl "http://localhost:8080/api/schedules?origin=Jakarta"
```

### 2.2 View seat availability for a schedule

```
GET /api/schedules/{scheduleId}/seats
```

Response — `200 OK`:
```json
[
  { "id": 1, "schedule_id": 1, "seat_number": "A1", "status": "available" },
  { "id": 2, "schedule_id": 1, "seat_number": "A2", "status": "locked", "locked_by": 3, "lock_expires_at": "2026-08-01T08:35:00Z" },
  { "id": 3, "schedule_id": 1, "seat_number": "A3", "status": "booked", "locked_by": 1 }
]
```

`status` is one of `available`, `locked`, `booked`. Note that `locked_by` is
kept populated even after a seat becomes `booked` — it's repurposed from
"who currently holds the temporary lock" to "who owns the confirmed
booking". This is what lets the frontend tell "booked by me" (shown in a
distinct color) apart from "booked by someone else", by comparing
`locked_by` against the logged-in user's own ID (returned at login). Calling
this endpoint also lazily releases any locks that have already expired, so
it always reflects accurate real-time availability.

`curl` example:
```bash
curl http://localhost:8080/api/schedules/1/seats
```

---

## 3. Protected endpoints (require a valid access token)

### 3.1 Lock a seat (5-minute hold)

```
POST /api/seats/{seatId}/lock
Authorization: Bearer <access_token>
```

Success response — `200 OK`:
```json
{
  "seat_id": 1,
  "status": "locked",
  "lock_expires_at": "2026-08-01T08:35:00Z"
}
```

Failure — `409 Conflict` (seat already locked by someone else, or already booked):
```json
{ "error": "seat is currently locked by another user" }
```

Failure — `401 Unauthorized` (missing/invalid/expired token):
```json
{ "error": "missing or malformed authorization header" }
```

`curl` example:
```bash
curl -X POST http://localhost:8080/api/seats/1/lock \
  -H "Authorization: Bearer <paste access_token here>"
```

### 3.2 Cancel a pending hold (before confirming)

Lets a user voluntarily release their own lock — e.g. they clicked the
wrong seat, or want to pick a different one — without waiting for the
5-minute expiry. Only releases the seat if the caller is the one currently
holding the lock; has no effect on a seat locked by someone else or already
booked.

```
POST /api/seats/{seatId}/unlock
Authorization: Bearer <access_token>
```

Success response — `200 OK`:
```json
{ "seat_id": 1, "status": "released" }
```

If the lock had already expired or wasn't yours, this is a harmless no-op
(still `200 OK`):
```json
{ "seat_id": 1, "status": "released", "note": "seat was not held by you (already expired or released)" }
```

`curl` example:
```bash
curl -X POST http://localhost:8080/api/seats/1/unlock \
  -H "Authorization: Bearer <paste access_token here>"
```

### 3.3 Confirm a booking

Must be called by the same user who holds a still-valid lock on the seat.

```
POST /api/seats/{seatId}/confirm
Authorization: Bearer <access_token>
```

Success response — `200 OK`:
```json
{
  "booking_id": 7,
  "seat_id": 1,
  "schedule_id": 1,
  "status": "confirmed"
}
```

Failure — `409 Conflict` (lock expired, seat not locked by you, or already booked):
```json
{ "error": "seat lock has expired; please lock the seat again" }
```

`curl` example:
```bash
curl -X POST http://localhost:8080/api/seats/1/confirm \
  -H "Authorization: Bearer <paste access_token here>"
```

---

## 4. Full example flow (copy-paste in order)

```bash
BASE=http://localhost:8080

# 1. Log in and capture the access token
TOKEN=$(curl -s -X POST $BASE/api/login \
  -d '{"username":"alice","password":"password123"}' \
  | python3 -c "import sys,json;print(json.load(sys.stdin)['access_token'])")

# 2. Search for schedules
curl -s "$BASE/api/schedules?origin=Jakarta"

# 3. View seats for schedule 1
curl -s "$BASE/api/schedules/1/seats"

# 4. Lock seat 1
curl -s -X POST "$BASE/api/seats/1/lock" -H "Authorization: Bearer $TOKEN"

# 5. Confirm the booking within 5 minutes
curl -s -X POST "$BASE/api/seats/1/confirm" -H "Authorization: Bearer $TOKEN"
```

---

## 5. Error format

All errors follow the same shape:
```json
{ "error": "human-readable message" }
```

| Status | Meaning |
|--------|---------|
| 400 | Malformed request body or query params |
| 401 | Missing, invalid, or expired token; wrong login credentials |
| 404 | Schedule or seat not found |
| 409 | Seat already locked/booked, or your lock has expired |
| 429 | Rate limit exceeded on `/api/schedules` |
| 500 | Unexpected server error |
