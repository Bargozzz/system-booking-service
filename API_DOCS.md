# API Documentation — Mini Booking Service

Base URL: `http://localhost:8000/api`

---

## 1. Authentication

### `POST /login`
Membuka sesi autentikasi dan menghasilkan token JWT Bearer.

- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "email": "user1@example.com",
    "password": "password"
  }
  ```
- **Response `200 OK`:**
  ```json
  {
    "message": "Login berhasil",
    "token": "eyJhbGciOiJIUzI1Ni...",
    "token_type": "bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "name": "User Satu",
      "email": "user1@example.com"
    }
  }
  ```

---

## 2. Schedules & Seat Map

### `GET /schedules/search`
Mencari jadwal bus berdasarkan kota asal, tujuan, dan tanggal.

- **Auth Required:** No (Throttled `throttle:search`)
- **Query Parameters:**
  - `origin` (string, required): e.g. `Jakarta`
  - `destination` (string, required): e.g. `Bandung`
  - `date` (date string YYYY-MM-DD, required): e.g. `2026-08-01`
- **Response `200 OK`:**
  ```json
  {
    "data": [
      {
        "id": 1,
        "operator": "Tiketux Express",
        "operator_code": "TKX",
        "origin": "Jakarta",
        "destination": "Bandung",
        "departure_date": "2026-08-01",
        "departure_time": "08:00:00",
        "price": 150000,
        "available_seats": 20,
        "total_seats": 20
      }
    ]
  }
  ```

---

### `GET /schedules/{id}/seats`
Mengambil layout dan status real-time 20 kursi (A1-E4) pada suatu jadwal.

- **Auth Required:** No
- **Response `200 OK`:**
  ```json
  {
    "schedule": {
      "id": 1,
      "operator": "Tiketux Express",
      "origin": "Jakarta",
      "destination": "Bandung",
      "price": 150000
    },
    "seats": [
      {
        "id": 1,
        "seat_number": "A1",
        "status": "available",
        "locked_until": null
      },
      {
        "id": 2,
        "seat_number": "A2",
        "status": "locked",
        "locked_until": "2026-07-30T01:30:00.000000Z"
      }
    ]
  }
  ```

---

## 3. Seat Locking & Booking Confirmation

### `POST /seats/{id}/lock`
Mengunci kursi secara atomik selama 5 menit untuk pengguna yang sedang terautentikasi.

- **Auth Required:** Yes (`Bearer <token>`)
- **Response `200 OK`:**
  ```json
  {
    "message": "Kursi berhasil dikunci selama 5 menit",
    "seat": {
      "id": 1,
      "seat_number": "A1",
      "status": "locked",
      "locked_until": "2026-07-30T01:30:00.000000Z"
    },
    "schedule_id": 1
  }
  ```
- **Response `409 Conflict`:**
  ```json
  {
    "message": "Kursi tidak tersedia. Sudah dikunci atau dipesan oleh pengguna lain."
  }
  ```

---

### `POST /seats/{id}/confirm`
Mengonfirmasi pemesanan kursi yang sedang dalam status dikunci oleh pengguna.

- **Auth Required:** Yes (`Bearer <token>`)
- **Response `200 OK`:**
  ```json
  {
    "message": "Pemesanan berhasil dikonfirmasi!",
    "booking": {
      "id": 1,
      "booking_code": "BK-X9Y2A1B8",
      "status": "confirmed",
      "seat_number": "A1",
      "operator": "Tiketux Express",
      "origin": "Jakarta",
      "destination": "Bandung",
      "departure_date": "2026-08-01",
      "departure_time": "08:00:00",
      "price": 150000,
      "created_at": "2026-07-30T01:26:00.000000Z"
    }
  }
  ```
- **Response `410 Gone`:**
  ```json
  {
    "message": "Waktu kunci kursi telah habis. Silakan coba lagi."
  }
  ```
