# Customer Order Web App

Aplikasi web-based untuk login JWT + TOTP 2FA dan CRUD data customer, order, serta order tracking.

## Stack

- Backend: Go, Gin, Gorm, SQLite, JWT, TOTP, godotenv
- Frontend: Vue 3, Vite, Naive UI, Pinia, Vue Router, Axios

## Menjalankan Backend

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

Backend berjalan di `http://localhost:8080`.

## Menjalankan Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Frontend berjalan di `http://localhost:5173`.

## Build Single Executable

Frontend production build bisa digabungkan ke binary backend Go melalui `embed`.

```bash
make build
```

Perintah tersebut menjalankan:

- `npm run build` di folder `frontend`
- output Vite ke `backend/internal/web/dist`
- `go build` backend ke `bin/customer-order-app`

Jalankan executable:

```bash
./bin/customer-order-app
```

Setelah itu aplikasi frontend dan API tersedia dari server yang sama:

- Web app: `http://localhost:8080`
- API: `http://localhost:8080/api`

Untuk development, `npm run dev` tetap memakai Vite di `http://localhost:5173` dan proxy `/api` ke backend `http://localhost:8080`.

## Login Development

- Email: `admin@example.com`
- Password: `admin12345`
- TOTP secret: `JBSWY3DPEHPK3PXP`

Masukkan secret tersebut ke Google Authenticator, Microsoft Authenticator, 1Password, Bitwarden, atau aplikasi TOTP lain. Login membutuhkan kode 6 digit yang berubah setiap 30 detik.

## Endpoint Utama

- `POST /api/auth/login`
- `GET /api/auth/me`
- `GET|POST /api/customers`
- `GET|PUT|DELETE /api/customers/:id`
- `GET|POST /api/orders`
- `GET|PUT|DELETE /api/orders/:id`
- `GET|POST /api/order-tracks`
- `GET|PUT|DELETE /api/order-tracks/:id`

Semua endpoint selain health check dan login membutuhkan header:

```http
Authorization: Bearer <jwt-token>
```

## Konfigurasi Penting

Backend membaca konfigurasi dari `backend/.env`:

- `APP_PORT`
- `APP_FRONTEND_ORIGIN`
- `JWT_SECRET`
- `JWT_EXPIRES_MINUTES`
- `DATABASE_DSN`
- `ADMIN_EMAIL`
- `ADMIN_PASSWORD`
- `ADMIN_TOTP_SECRET`

Untuk production, ganti `JWT_SECRET`, password admin, dan TOTP secret. SQLite dipakai agar development cepat; Gorm membuat migrasi otomatis saat server start.
