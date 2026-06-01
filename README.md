# Customer Order Web App

Aplikasi web-based untuk login JWT + TOTP 2FA dan CRUD data customer, order, serta order tracking.

## Stack

- Backend dan UI server-rendered: Go, Gin, Gorm, SQLite, html/template, JWT, TOTP, godotenv
- UI interactivity: HTMX dan AlpineJS
- CSS framework: Tailwind CSS
- UI kit: Preline UI

## Menjalankan Aplikasi

```bash
cp backend/.env.example backend/.env
make run
```

Aplikasi berjalan di `http://localhost:8080`.

## Build Executable

```bash
make build
./bin/customer-order-app
```

Tidak ada proses build JavaScript. Go meng-embed template HTML dari `backend/internal/web/templates`, lalu executable menyajikan halaman server-rendered di origin yang sama dengan API.

## Login Development

- Email: `admin@example.com`
- Password: `admin12345`
- TOTP secret: `JBSWY3DPEHPK3PXP`

Masukkan secret tersebut ke Google Authenticator, Microsoft Authenticator, 1Password, Bitwarden, atau aplikasi TOTP lain. Login membutuhkan kode 6 digit yang berubah setiap 30 detik.

## Endpoint Utama

- `GET /login`
- `POST /login`
- `GET /dashboard`
- `GET|POST /customers`
- `GET /customers/:id/edit`
- `POST /customers/:id/update`
- `DELETE /customers/:id`
- `GET|POST /orders`
- `GET /orders/:id/edit`
- `POST /orders/:id/update`
- `DELETE /orders/:id`
- `GET|POST /tracks`
- `GET /tracks/:id/edit`
- `POST /tracks/:id/update`
- `DELETE /tracks/:id`
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
