# Mini Booking Service

A simplified booking application (inspired by Tiketux) demonstrating schedule
search, temporary seat locking, and race-condition-safe booking confirmation.

- **Backend**: Go (standard library `net/http`, SQLite via `mattn/go-sqlite3`, JWT via `golang-jwt/jwt`)
- **Frontend**: Vue 3 + Vite, Vue Router, Axios

📄 Full API reference with request/response examples: **[API.md](./API.md)**
📮 Postman collection: **[postman_collection.json](./postman_collection.json)**

---

## 1. Project layout

```
mini-booking/
├── docker-compose.yml    runs backend + frontend together
├── API.md                 API reference with curl examples
├── postman_collection.json
├── scripts/
│   └── loadtest.go         fires N concurrent lock requests at 1 seat
├── backend/               Go REST API
│   ├── main.go              entrypoint + route table
│   ├── main_test.go          integration tests (httptest)
│   ├── Dockerfile
│   ├── db/                   sqlite connection, migrations, seed data
│   ├── models/                data structs
│   ├── middleware/             JWT auth
│   ├── ratelimit/               per-IP token-bucket rate limiter
│   └── handlers/                 auth, schedules, bookings, lock-expiry sweeper
└── frontend/               Vue 3 SPA
    ├── Dockerfile
    ├── nginx.conf
    └── src/
        ├── api/               axios client (auto token refresh)
        ├── store/              auth session state
        ├── router/              vue-router + auth guard
        └── views/                Login, Search, SeatSelection, BookingSummary, Confirmation
```

---

## 2. Running the backend (local, without Docker)

Requires Go 1.22+ and a C compiler (for the sqlite3 driver's cgo build).

```bash
cd backend
go run .
```

The server starts on **http://localhost:8080**, creates `minibooking.db` in the
current directory on first run, and seeds it with:

- Two dummy users: `alice` / `bob`, both with password `password123`
- One operator ("Tiketux Air"), two routes (Jakarta→Bandung, Bandung→Surabaya)
- Three schedules across those routes, each with 12 seats (A1–A6, B1–B6)

To reset the data, stop the server and delete `minibooking.db*`.

### Running the tests

```bash
cd backend
go test ./... -v          # functional tests
go test ./... -race       # same tests, with Go's data-race detector
```

Covers: successful end-to-end booking, lock conflicts between two users,
concurrent-lock race safety (15 simultaneous lock attempts → exactly 1
succeeds), missing/invalid JWTs, automatic lock expiry, and bad login
credentials.

---

## 3. Running the frontend (local, without Docker)

Requires Node.js 18+.

```bash
cd frontend
npm install
npm run dev
```

Opens on **http://localhost:5173** (or similar) and talks to the backend at
`http://localhost:8080` — make sure the backend is running first.

For a production build:

```bash
npm run build
npm run preview
```

---

## 4. Running everything with Docker (bonus)

Requires Docker and Docker Compose. This is the fastest way to run the whole
stack with one command — no need to install Go, Node.js, or a C compiler
locally.

```bash
docker compose up --build
```

- Backend: **http://localhost:8080**
- Frontend: **http://localhost:5173**

The SQLite database persists across restarts in a named Docker volume
(`booking-data`). To wipe it and reseed from scratch:

```bash
docker compose down -v
docker compose up --build
```

To stop everything:

```bash
docker compose down
```

---

## 5. Load testing concurrency (bonus)

`scripts/loadtest.go` fires N concurrent `POST /api/seats/{id}/lock` requests
at a single seat and asserts that exactly one succeeds — a direct,
observable proof that the race-condition handling works, not just a claim.

With the backend running (locally or via Docker) on port 8080:

```bash
cd scripts
go run loadtest.go                    # default: 50 concurrent requests
go run loadtest.go -requests=100      # customize the count
go run loadtest.go -seat=5            # target a specific seat ID
```

Example output:

```
== Mini Booking Service — concurrency load test ==
Target: http://localhost:8080 | concurrent requests: 50

✓ Logged in as alice
✓ Targeting seat ID 1

Firing 50 concurrent POST /api/seats/1/lock requests...

Done in 26.17ms.

Result breakdown:
  200 OK (lock acquired)                        x1
  409 Conflict (seat already locked/booked)     x49

✅ PASS: exactly 1 request acquired the lock, as expected.
```

---

## 6. API reference

See **[API.md](./API.md)** for the full reference with example requests and
responses for every endpoint, or import **[postman_collection.json](./postman_collection.json)**
into Postman (it auto-saves the token from the Login request into a
collection variable so every protected request picks it up automatically).

Quick summary:

| Method | Path                          | Auth | Description                                  |
|--------|-------------------------------|------|------------------------------------------------|
| POST   | `/api/login`                  | none | Returns `access_token` (15 min) + `refresh_token` (24 h) |
| POST   | `/api/refresh`                | none | Exchanges a valid refresh token for a new access token |
| GET    | `/api/schedules`               | none | Search by `?origin=&destination=&date=YYYY-MM-DD`, rate-limited to 2 req/s (burst 10) per IP |
| GET    | `/api/schedules/{id}/seats`     | none | Seat map for a schedule (also lazily releases expired locks) |
| POST   | `/api/seats/{id}/lock`         | JWT  | Locks a seat for 5 minutes if available or its previous lock expired |
| POST   | `/api/seats/{id}/confirm`      | JWT  | Confirms a booking; must hold a valid, unexpired lock on that seat |

---

## 7. Concurrency approach

Seat locking and booking use **atomic conditional `UPDATE ... WHERE`**
statements rather than read-then-write logic. For example, locking a seat
runs:

```sql
UPDATE seats
SET status = 'locked', locked_by = ?, lock_expires_at = ?
WHERE id = ? AND (status = 'available' OR (status = 'locked' AND lock_expires_at < ?))
```

Whichever concurrent request's `UPDATE` commits first flips the row's state;
every other request's `WHERE` clause then no longer matches, so its
`RowsAffected` is 0 and it's rejected with `409 Conflict`. This holds
regardless of how many app-server instances are running, since the guarantee
comes from the database, not from an in-process mutex — an important
distinction, because an in-process mutex would *not* protect against two
separate backend instances (e.g. two Docker replicas) racing for the same
seat, while a database-level atomic update does.

Booking confirmation follows the same pattern inside a transaction: it
re-checks the lock owner and expiry, then atomically flips the seat to
`booked` and inserts the booking row in one commit, so the two can never
drift out of sync (e.g. a seat marked `booked` with no matching booking
record, or vice versa).

A background sweeper also runs every 30 seconds to proactively release any
expired locks, and `GET /api/schedules/{id}/seats` lazily releases expired
locks on read too, so availability is always accurate even between sweeps.

This was verified two ways:
- **Automated test** (`TestConcurrentLockingIsRaceFree` in `main_test.go`):
  15 goroutines lock the same seat simultaneously; exactly 1 succeeds.
- **Load test script** (`scripts/loadtest.go`, see section 5): 50 real HTTP
  requests fired concurrently at one seat; exactly 1 succeeds.

---

## 8. Authentication approach

- **JWT, stateless, no session store.** Each token is a signed HS256 JWT
  carrying the user ID, username, and a `type` claim (`access` or
  `refresh`), so the server never needs to look up sessions in the database
  to validate a request — it only verifies the signature and expiry.
- **Pre-seeded dummy users, no registration**, per the brief. Passwords are
  stored in plaintext in this demo (`db/db.go`); a production system would
  hash them (bcrypt/argon2) instead.
- **Short-lived access token (15 min) + longer-lived refresh token (24 h).**
  This limits how long a stolen access token stays useful, while the refresh
  token lets the frontend stay logged in for a workday without repeated
  logins. The `type` claim prevents a refresh token from being used directly
  as an access token even if someone tries to pass it to a protected
  endpoint.
- **Selective protection.** `RequireAuth` middleware (`middleware/auth.go`)
  wraps only the seat-lock and booking-confirmation handlers, matching the
  brief's requirement that search and seat-availability stay public.
- **Frontend auto-refresh.** The axios client (`frontend/src/api/client.js`)
  intercepts `401` responses, silently exchanges the refresh token for a new
  access token, retries the original request once, and only forces the user
  back to the login screen if the refresh token itself is invalid/expired.

---

## 9. Publishing to GitHub

This project isn't in a Git repository yet. To push it as your own
repository:

```bash
cd mini-booking          # the folder you extracted from the zip
git init
git add .
git commit -m "Initial commit: Mini Booking Service"
```

Then create an empty repository on GitHub (via the website: click **New**
at github.com/new — do **not** initialize it with a README), and push:

```bash
git remote add origin https://github.com/<your-username>/<repo-name>.git
git branch -M main
git push -u origin main
```

`node_modules/`, build output (`dist/`), and the local SQLite database files
are already excluded via `.gitignore` in both `backend/` and `frontend/`, so
they won't be committed.

---

## 10. Notes on scope / simplifications

- Authentication uses two pre-seeded dummy users with plaintext passwords, as
  specified in the brief — this is a demo, not a production auth system.
- SQLite runs with a single-connection pool and `_busy_timeout`, which is
  enough for a demo's concurrency needs; a high-throughput production system
  would use Postgres with `SELECT ... FOR UPDATE` or advisory locks instead.
- Rate limiting is a simple in-memory per-IP token bucket; it resets if the
  server restarts and doesn't share state across multiple instances (a
  production deployment with multiple replicas would use a shared store like
  Redis instead).
- Stack choice: the brief allows "stack bebas" (any stack); this project uses
  Go + Vue instead of the suggested Node.js/PHP, since Go's standard library
  and static typing made the concurrency guarantees easier to demonstrate and
  test explicitly.
