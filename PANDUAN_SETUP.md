# Panduan Setup — Mini Booking Service (Windows)

Panduan ini untuk Anda yang baru pertama kali menjalankan project seperti ini.
Ikuti dari atas ke bawah, **jangan skip langkah**, dan jangan panik kalau ada
error — di bagian akhir ada tabel troubleshooting.

---

## Yang Anda butuhkan sebelum mulai

| Software | Kenapa dibutuhkan | Link download |
|---|---|---|
| Go 1.22+ | Menjalankan backend (server) | https://go.dev/dl/ |
| Node.js LTS | Menjalankan frontend (tampilan web) | https://nodejs.org |
| w64devkit (compiler C) | Dibutuhkan database SQLite di backend | https://github.com/skeeto/w64devkit/releases |
| VSCode | Text editor untuk membuka & menjalankan code | https://code.visualstudio.com |

---

## Tahap 1 — Ekstrak ZIP

1. Cari file `mini-booking-service.zip` yang sudah di-download (biasanya di folder **Downloads**)
2. Klik kanan file tersebut → **Extract All...**
3. Pilih lokasi tujuan, misalnya `D:\KAMPUS\tutor\mini-booking-services`
4. Klik **Extract**

Setelah ini Anda akan punya folder berisi: `backend/`, `frontend/`, `scripts/`, `README.md`, `API.md`, `docker-compose.yml`.

---

## Tahap 2 — Install Go

1. Buka https://go.dev/dl/
2. Download `go1.xx.x.windows-amd64.msi`
3. Jalankan installer, klik **Next** terus sampai **Finish** (pengaturan default sudah benar)
4. **Restart VSCode** (atau restart komputer) setelah ini — supaya perintah `go` bisa dikenali di terminal

**Verifikasi**: buka terminal baru di VSCode, ketik:
```powershell
go version
```
Harus muncul teks seperti `go version go1.22.x windows/amd64`.

---

## Tahap 3 — Install Node.js

1. Buka https://nodejs.org, download versi **LTS**
2. Jalankan installer, klik **Next** terus sampai selesai

**Verifikasi**:
```powershell
node --version
npm --version
```

---

## Tahap 4 — Install compiler C (w64devkit)

Backend memakai database SQLite, dan driver SQLite di Go butuh compiler C untuk di-build. Ini langkah yang sering terlewat pemula, jadi ikuti persis:

1. Buka https://github.com/skeeto/w64devkit/releases
2. Download file `w64devkit-x.x.x.zip` (ukuran sekitar 100–200MB)
3. Extract ke `C:\w64devkit` (klik kanan zip → Extract All → arahkan ke `C:\`)
4. Pastikan setelah extract ada file `C:\w64devkit\bin\gcc.exe`
5. Tambahkan ke PATH:
   - Tekan tombol **Windows**, ketik `env`, pilih **"Edit the system environment variables"**
   - Klik **Environment Variables...**
   - Cari baris **Path** (di User variables atau System variables) → klik **Edit** → **New**
   - Ketik: `C:\w64devkit\bin`
   - Klik **OK** di semua jendela (bukan Cancel)
6. **Tutup VSCode sepenuhnya dan buka ulang** (PATH baru tidak terbaca oleh terminal yang sudah terbuka)

**Verifikasi**:
```powershell
gcc --version
```

---

## Tahap 5 — Buka project di VSCode

1. Buka VSCode
2. **File → Open Folder...** → pilih folder hasil extract, misalnya `D:\KAMPUS\tutor\mini-booking-services`
3. (Opsional) Install extension **Go** (dari tim golang.org) dan **Volar** (untuk Vue) lewat tab Extensions

---

## Tahap 6 — Jalankan Backend

Buka terminal baru di VSCode (**Terminal → New Terminal**):

```powershell
cd backend
go mod download
go run .
```

Kalau berhasil, akan muncul log seperti ini dan **terminal ini harus tetap terbuka**:

```
Seeding database with dummy data...
Database initialized and seeded.
Lock expiry sweeper started (interval: 30s).
Mini Booking Service listening on :8080
```

Ini artinya backend jalan di `http://localhost:8080`.

---

## Tahap 7 — Jalankan Frontend

Backend harus tetap jalan, jadi buka **terminal kedua** (klik ikon **+** di panel terminal, jangan ganti terminal yang lama):

```powershell
cd frontend
npm install
npm run dev
```

Tunggu sampai muncul:
```
➜  Local:   http://localhost:5173/
```

---

## Tahap 8 — Buka di browser

1. Buka Chrome/Edge
2. Kunjungi `http://localhost:5173`
3. Login dengan:
   - Username: `alice`
   - Password: `password123`
4. Coba alur: cari jadwal → pilih kursi → konfirmasi booking

---

## Ringkasan: jalankan ulang di lain waktu

Setiap kali mau menjalankan lagi (setelah setup awal selesai):

- **Terminal 1**: `cd backend` → `go run .`
- **Terminal 2**: `cd frontend` → `npm run dev`
- Buka `http://localhost:5173`

---

## (Opsional) Menjalankan dengan Docker

Kalau Anda sudah install Docker Desktop, ini cara paling singkat — tidak perlu install Go/Node/gcc sama sekali:

```powershell
docker compose up --build
```

Backend otomatis jalan di `:8080`, frontend di `:5173`.

---

## (Opsional) Membuktikan sendiri penanganan race condition

Ini skrip yang menembak 50 request bersamaan ke 1 kursi, untuk membuktikan hanya 1 yang berhasil (berguna untuk didemokan saat interview):

```powershell
cd scripts
go run loadtest.go
```

---

## Troubleshooting

| Error | Penyebab | Solusi |
|---|---|---|
| `'go' is not recognized` | Go belum terpasang / PATH belum terbaca | Install ulang Go, restart VSCode |
| `cgo: C compiler "gcc" not found` | w64devkit belum terpasang / PATH belum benar | Ulangi Tahap 4, pastikan `C:\w64devkit\bin\gcc.exe` ada, restart VSCode total |
| `Cannot find native binding` saat `npm run dev` | Bug npm terkait optional dependencies | `Remove-Item -Recurse -Force node_modules`, `Remove-Item package-lock.json`, lalu `npm install` ulang |
| `port 8080 already in use` | Ada proses lama masih jalan | Tutup semua terminal yang menjalankan backend, coba lagi |
| Halaman browser blank / gagal konek API | Backend tidak jalan | Cek terminal 1, pastikan tidak ada error, dan masih menampilkan "listening on :8080" |
| `npm install` gagal karena koneksi | Masalah jaringan | Cek internet, coba `npm cache clean --force` lalu ulangi |

Kalau masih ada error yang belum tercantum di sini, salin pesan error persis dan tanyakan — jangan coba tebak-tebak sendiri dulu.
