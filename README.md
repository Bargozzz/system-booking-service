# Mini Booking Service

A simplified booking application (inspired by Tiketux) demonstrating schedule
search, temporary seat locking, and race-condition-safe booking confirmation.

- **Backend**: Go (standard library `net/http`, SQLite via `mattn/go-sqlite3`, JWT via `golang-jwt/jwt`)
- **Frontend**: Vue 3 + Vite, Vue Router, Axios

📄 Full API reference with request/response examples: **[API.md](./API.md)**
📮 Postman collection: **[postman_collection.json](./postman_collection.json)**

---

## Quick Start (recommended for reviewers)

The fastest, most reliable way to run this project on any machine —
Windows, macOS, or Linux — is Docker, because it avoids every
OS-specific setup issue described further below (missing C compiler, npm
platform-binary bugs, etc.):

```bash
git clone <this-repo-url>
cd mini-booking-service
docker compose up --build
```

Wait for both services to report they're up (first build takes a few
minutes while it downloads base images), then open:

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Login with**: `alice` / `password123` (or `bob` / `password123`)

If you don't have Docker installed, or prefer running things natively,
follow sections 2 and 3 below instead — they cover setup per operating
system, with a troubleshooting section for the specific errors that come
up most often.

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

### Prerequisites

| Requirement | Why | How to check | How to install |
|---|---|---|---|
| **Go 1.22+** | Compiles and runs the server | `go version` | https://go.dev/dl/ |
| **A C compiler (cgo)** | The SQLite driver (`mattn/go-sqlite3`) needs to compile a small amount of C code | `gcc --version` | See OS-specific instructions below |

**Installing a C compiler:**
- **Windows**: Download [w64devkit](https://github.com/skeeto/w64devkit/releases) (a portable
  GCC toolchain — no installer needed), extract it anywhere (e.g. `C:\w64devkit`),
  then add `C:\w64devkit\bin` to your `PATH` (Windows key → search "environment
  variables" → Edit the system environment variables → Environment Variables →
  edit `Path` → add the folder). **Restart your terminal/IDE** after this so the
  updated `PATH` is picked up.
- **macOS**: Install Xcode Command Line Tools: `xcode-select --install`
- **Linux (Debian/Ubuntu)**: `sudo apt-get install build-essential`

### Steps

```bash
cd backend
go mod download
go run .
```

Expected output:
```
2026/08/01 10:00:00 Seeding database with dummy data...
2026/08/01 10:00:00 Database initialized and seeded.
2026/08/01 10:00:00 Lock expiry sweeper started (interval: 30s).
2026/08/01 10:00:00 Mini Booking Service listening on :8080
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
succeeds), missing/invalid JWTs, automatic lock expiry, voluntarily
cancelling a hold, and bad login credentials. All 7 should pass.

### Backend troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `cgo: C compiler "gcc" not found` | No C compiler installed/on PATH | See "Installing a C compiler" above; **restart your terminal** after adding it to PATH |
| `bind: address already in use` (port 8080) | Another process is already using port 8080 | Stop the other process, or run with a different port: `PORT=8081 go run .` (macOS/Linux) or `$env:PORT="8081"; go run .` (PowerShell) |
| `database is locked` | Two processes (e.g. this server + a SQLite GUI tool) writing at the same time | Close the other tool, or wait — the connection is configured with `_busy_timeout=5000` so it retries for 5s before failing |

---

## 3. Running the frontend (local, without Docker)

### Prerequisites

| Requirement | Why | How to check | How to install |
|---|---|---|---|
| **Node.js 18+** | Runs Vite (dev server + build tool) | `node --version` | https://nodejs.org (LTS version) |

The backend must already be running (see section 2) before the frontend
can successfully load data — start it first if you haven't.

### Steps

```bash
cd frontend
npm install
npm run dev
```

Expected output:
```
  VITE vX.X.X  ready in XXX ms

  ➜  Local:   http://localhost:5173/
```

Open **http://localhost:5173** in a browser. It talks to the backend at
`http://localhost:8080` by default (configurable via `VITE_API_BASE_URL`,
see `.env.example`).

For a production build:

```bash
npm run build
npm run preview
```

### Frontend troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `Cannot find module @rollup/rollup-<platform>` or `Cannot find native binding` (mentions Rolldown/Rollup) | Known npm bug with optional platform-specific dependencies ([npm/cli#4828](https://github.com/npm/cli/issues/4828)) — happens more often on Windows | Delete `node_modules` and `package-lock.json`, then run `npm install` again:<br>`Remove-Item -Recurse -Force node_modules` (PowerShell) or `rm -rf node_modules` (macOS/Linux), then `Remove-Item package-lock.json` / `rm package-lock.json`, then `npm install` |
| Same error persists after the above | The installed Vite version's bundled Rolldown binary genuinely isn't available for your platform yet | Pin an older, stable Vite version: `npm uninstall vite && npm install vite@^5.4.0 --save-dev`, then `npm run dev` again |
| Page loads but shows no schedules / network errors in browser console | Backend isn't running, or is running on a different port than `VITE_API_BASE_URL` expects | Confirm backend is up at http://localhost:8080 (visit it directly, or `curl http://localhost:8080/api/schedules`) |
| `EBADENGINE` warning when running `npm install -g npm@latest` | Your Node.js version is older than what the newest npm requires | Safe to ignore — the npm version that ships with your Node install is sufficient for this project; there's no need to upgrade npm globally |

---

## 4. Running everything with Docker

### Prerequisites

| Requirement | How to check | How to install |
|---|---|---|
| **Docker Desktop** (includes Docker Compose) | `docker --version` and `docker compose version` | https://www.docker.com/products/docker-desktop/ |

On Windows, Docker Desktop requires WSL2 — the installer will prompt you
to enable/install it if needed. **Docker Desktop must be running** (check
for its whale icon in the system tray/menu bar) before using `docker` commands.

### Steps

```bash
docker compose up --build
```

First build takes a few minutes (downloading Go, Node, and nginx base
images). Subsequent runs are much faster thanks to Docker's layer cache.

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

### Docker troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `Cannot connect to the Docker daemon` | Docker Desktop isn't running | Open the Docker Desktop application and wait for it to fully start |
| Frontend build fails with a Rollup/native-binding error (same family as the npm issue above) | A `node_modules` folder from your host machine (built for your OS) got copied into the Linux container, overwriting the correct Linux binaries | Make sure `frontend/.dockerignore` exists and contains `node_modules` (it's included in this repo already); then clear the build cache and rebuild: `docker compose down -v && docker builder prune -f && docker compose up --build` |
| `port is already allocated` | Something else on your machine is already using port 8080 or 5173 (e.g. you also have the backend/frontend running natively) | Stop the natively-running processes first, or edit the port mappings in `docker-compose.yml` |
| Old code/data still showing after a rebuild | Docker reused cached layers | `docker compose down -v` (removes the data volume too) then `docker compose up --build` |

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
| POST   | `/api/seats/{id}/unlock`       | JWT  | Voluntarily releases the caller's own pending hold before confirming (used by the frontend's "Cancel" button) |

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
