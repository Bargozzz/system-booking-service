# Product Requirement Document (PRD): Mini Booking Service

## 1. Overview & Objective
Mini Booking Service is a full-stack bus ticket reservation system designed to handle real-time schedule searches, seat selection, and high-concurrency seat reservation with temporary locking.

---

## 2. Key Features & Requirements

### 2.1 Authentication & Authorization
- User login using Email and Password (`POST /api/login`).
- Session management using Bearer JWT / API Token.
- Protected endpoints (`/me`, `/logout`, `/seats/{id}/lock`, `/seats/{id}/confirm`).

### 2.2 Schedule Search & Seat Availability
- Search bus routes by **Origin**, **Destination**, and **Departure Date** (`GET /api/schedules/search`).
- Rate limiting applied to search endpoint (`throttle:search`).
- View seat layout for specific schedule (`GET /api/schedules/{id}/seats`).

### 2.3 Atomic Seat Locking (Concurrency Handling)
- Lock seat temporarily for **5 minutes** (`POST /api/seats/{id}/lock`).
- Utilizes database transactions (`DB::transaction`) and pessimistic locking (`lockForUpdate`) to prevent race conditions.
- Auto-expiration mechanism for locks older than 5 minutes.

### 2.4 Booking Confirmation
- User can confirm booking for locked seat (`POST /api/seats/{id}/confirm`).
- Generates unique booking code (`BK-XXXXXXXX`).
- Updates seat status to `booked`.

---

## 3. Database Schema

- **`operators`**: `id`, `name`, `code`, `timestamps`
- **`bus_routes`**: `id`, `operator_id`, `origin`, `destination`, `price`, `timestamps`
- **`schedules`**: `id`, `bus_route_id`, `departure_date`, `departure_time`, `timestamps`
- **`seats`**: `id`, `schedule_id`, `seat_number`, `status` (`available`, `locked`, `booked`), `locked_by`, `locked_until`, `timestamps`
- **`bookings`**: `id`, `user_id`, `seat_id`, `schedule_id`, `booking_code`, `status`, `timestamps`
- **`users`**: `id`, `name`, `email`, `password`, `timestamps`

---

## 4. API Endpoints Specification

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/login` | Authenticate user & get JWT token | No |
| `GET` | `/api/schedules/search` | Search schedules with origin, dest, date | No (Throttled) |
| `GET` | `/api/schedules/{id}/seats` | Get seat grid and current locking statuses | No |
| `GET` | `/api/me` | Fetch authenticated user details | Yes |
| `POST` | `/api/logout` | Revoke current token | Yes |
| `POST` | `/api/seats/{id}/lock` | Lock seat for 5 minutes | Yes |
| `POST` | `/api/seats/{id}/confirm` | Confirm booking of locked seat | Yes |

---

## 5. User Interface & Experience (Frontend Roadmap)

1. **Authentication Bar / Modal**: Easy login with pre-configured demo credentials.
2. **Search Engine Banner**: Filter routes by origin (e.g. Jakarta, Bandung, Surabaya, Malang), destination, and date.
3. **Schedule Cards**: Display operator badge, departure time, price, and remaining available seat count.
4. **Interactive Seat Matrix**:
   - Visual 4-column seat map (A1-E4).
   - Live state badges: Available (Green), Selected/Locked by User (Yellow + 5-min timer), Locked by Other/Booked (Gray/Disabled).
5. **Booking Receipt Modal**: Shows confirmed booking ticket with unique code, operator details, seat number, and departure info.
