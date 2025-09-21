# Backend Sewa dan Jual Kendaraan 🚗💨

**Backend API untuk Platform Penyewaan dan Jual-Beli Kendaraan di Sulawesi Tenggara**

API ini adalah tulang punggung untuk aplikasi marketplace otomotif multifungsi. Dibangun dengan Go (Golang) dan mengadopsi arsitektur berlapis (layered architecture) yang bersih, API ini dirancang untuk performa tinggi, skalabilitas, dan kemudahan pemeliharaan.

---

## ✨ Fitur Utama

Sistem ini memiliki fungsionalitas yang kaya, mencakup semua kebutuhan inti dari platform marketplace modern.

### 👤 **Autentikasi & Manajemen Pengguna**

- Registrasi pengguna dengan tiga peran berbeda: `customer`, `vendor`, dan `admin`.
- Sistem login aman menggunakan **JSON Web Tokens (JWT)**.
- Middleware untuk proteksi rute berdasarkan autentikasi dan peran (Role-Based Access Control).

### vehicle **Manajemen Listing & Pencarian**

- **CRUD** (Create, Read, Update, Delete) penuh untuk listing kendaraan oleh vendor.
- Upload gambar kendaraan yang terintegrasi langsung dengan **Cloudinary**.
- **Pencarian Lanjutan & Filter Dinamis** berdasarkan tipe, merek, transmisi, tahun, harga, dll.

### 📅 **Alur Kerja Penyewaan (Rental)**

- Sistem booking dengan pengecekan ketersediaan tanggal secara _real-time_ untuk mencegah tumpang tindih.
- Siklus hidup status booking yang lengkap (`pending_payment`, `confirmed`, `rented_out`, `completed`, `cancelled`).
- Simulasi integrasi _Payment Gateway_ dengan endpoint _callback_.
- Riwayat booking untuk customer dan vendor.

### 💸 **Alur Kerja Jual-Beli (Sales)**

- Fungsionalitas untuk memulai transaksi pembelian kendaraan.
- Perubahan status kendaraan menjadi `sold` setelah transaksi selesai, membuatnya tidak lagi tersedia di pasar.
- Riwayat transaksi penjualan untuk vendor dan pembelian untuk customer.

### ⭐ **Ulasan & Rating**

- Kemampuan bagi customer untuk memberikan rating (1-5 bintang) dan ulasan pada booking yang telah selesai.
- Endpoint publik untuk melihat semua ulasan dari sebuah kendaraan.

### 🛡️ **Panel Admin**

- Sistem **verifikasi vendor** oleh admin. Vendor yang belum terverifikasi tidak dapat memposting listing.
- Kemampuan admin untuk melihat dan menghapus pengguna (User Management).
- Kemampuan admin untuk melihat dan menghapus listing kendaraan (Listing Management).

### 💬 **Chat Real-time**

- Komunikasi dua arah secara instan antara pengguna menggunakan **WebSockets**.
- Sistem percakapan pribadi yang menyimpan riwayat pesan di database.
- REST API untuk memulai percakapan dan mengambil riwayat pesan.

---

## 🛠️ Tumpukan Teknologi

- **Bahasa:** Go (Golang)
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
- **Konektor DB:** pgx
- **Containerization:** Docker & Docker Compose
- **Autentikasi:** JSON Web Tokens (JWT)
- **Password Hashing:** Bcrypt
- **WebSockets:** Gorilla WebSocket
- **Penyimpanan Gambar:** Cloudinary
- **Manajemen Konfigurasi:** Environment Variables (.env)

---

## 🚀 Instalasi & Menjalankan Proyek

Untuk menjalankan proyek ini di lingkungan lokal, Anda memerlukan **Git** dan **Docker**.

**1. Clone Repositori**

```bash
git clone <URL_REPOSITORI_ANDA>
cd sultra-otomotif-api
```

**2. Konfigurasi Environment**

Buat file **.env** di direktori utama dengan menyalin dari contoh.

```bash
cp .env.example .env
```

**Kemudian, buka file .env dan isi semua variabel yang diperlukan:**

```bash
# Konfigurasi Database
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=sultra_otomotif
DB_PORT=5432
APP_PORT=8080
JWT_SECRET_KEY=iniadalahkuncirahasiayangSANGATpanjangdankuat
CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME
```

**3. Jalankan Aplikasi**
Gunakan Docker Compose untuk membangun dan menjalankan semua service (aplikasi Go & database Postgres).

```bash
docker-compose up --build
```

**Server API akan berjalan di http://localhost:8080.**

### 📁 Struktur Proyek

```bash
/
├── cmd/api/             # Entry point utama aplikasi (main.go)
├── internal/
│   ├── config/          # Manajemen konfigurasi (.env)
│   ├── handler/         # Layer untuk menangani HTTP request & response
│   ├── helper/          # Fungsi-fungsi bantuan (response, password, dll)
│   ├── middleware/      # Middleware (JWT Auth, Role Check)
│   ├── model/           # Definisi struct Go untuk data (User, Vehicle, dll)
│   ├── repository/      # Layer untuk interaksi langsung dengan database (SQL queries)
│   ├── service/         # Layer untuk logika bisnis utama
│   └── websocket/       # Logika untuk Hub dan Client WebSocket
├── .env                 # File konfigurasi (Jangan di-commit ke Git!)
├── .env.example         # Contoh file konfigurasi
├── go.mod               # Manajemen dependensi Go
├── Dockerfile           # Instruksi untuk membangun image aplikasi
└── docker-compose.yml   # Mendefinisikan dan menjalankan service app + db
```

### 📚 API Endpoints (Gambaran Umum)

Dokumentasi API lengkap dapat dibuat menggunakan Postman atau Swagger. Berikut adalah gambaran umum endpoint yang tersedia:

- **Auth:** /api/v1/auth/register, /api/v1/auth/login

- **Vehicles:** GET /vehicles, GET /vehicles/:id, POST /vehicles, PUT /vehicles/:id, DELETE /vehicles/:id, POST /vehicles/:id/images

- **Bookings:** POST /bookings, GET /bookings/my-bookings, GET /bookings/vendor, GET /bookings/:id, PATCH /bookings/:id/status

- **Sales:** POST /vehicles/:id/purchase, GET /sales/purchases, GET /sales/sales, POST /sales/callback

- **Reviews:** POST /bookings/:booking_id/reviews, GET /vehicles/:id/reviews

- **Admin:** GET /admin/vendors, PATCH /admin/vendors/:id/verify, GET /admin/users, DELETE /admin/users/:id, GET /admin/vehicles, DELETE /admin/vehicles/:id

- **WebSocket:** GET /api/v1/ws

`SELAMAT MENGGUNAKAN - SALAm HANGAT DARI SAYA`
