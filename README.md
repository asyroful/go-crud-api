# Go CRUD API (Gin & GORM)

API RESTful sederhana yang dibangun menggunakan Go dengan framework Gin. Proyek ini mengimplementasikan operasi CRUD (Create, Read, Update, Delete) dengan arsitektur bersih (Handlers, Services, Repository) dan menggunakan GORM sebagai ORM untuk berinteraksi dengan database PostgreSQL.

## Fitur

-   **Manajemen Pengguna**: Registrasi dan Login dengan sistem role (Admin/User).
-   **Role-Based Access Control**: 
    -   **Admin**: Akses penuh ke manajemen kategori, melihat semua transaksi, dan CRUD user management.
    -   **User**: Hanya bisa melihat kategori dan CRUD transaksi sendiri.
-   **Otentikasi JWT**: Endpoint diamankan menggunakan JSON Web Tokens.
-   **CRUD untuk Kategori**:
    -   Membuat, Membaca, Memperbarui, dan Menghapus kategori (Admin only).
    -   **Paginasi**: Mendukung paginasi (`limit` & `page`) untuk daftar kategori.
    -   **Pencarian**: Mendukung pencarian berdasarkan nama kategori (`q`).
-   **CRUD untuk Transaksi**:
    -   Membuat, Membaca, Memperbarui, dan Menghapus transaksi.
    -   **Filter Date Range**: Default filter tanggal 27 bulan lalu hingga 26 bulan ini.
    -   **Filter Tipe**: Filter berdasarkan tipe transaksi (income/expense).
    -   **Filter Kategori**: Filter berdasarkan kategori.
    -   Admin dapat melihat transaksi semua user.
-   **Balance/Saldo**: 
    -   Menampilkan total income, expense, dan balance berdasarkan range tanggal.
    -   Default range: 27 bulan lalu - 26 bulan ini.
-   **Admin User Management**: CRUD pengguna (khusus admin).
-   **Arsitektur Bersih**: Kode diorganisir ke dalam lapisan `handlers`, `services`, dan `repository`.
-   **Database PostgreSQL**: Menggunakan GORM untuk interaksi database.
-   **Manajemen Konfigurasi**: Menggunakan file `.env` untuk mengelola variabel lingkungan.
-   **Containerization**: Siap dijalankan menggunakan Docker dan Docker Compose.

## Struktur Proyek

```
.
├── config/
│   └── db.go           # Koneksi database
├── docs/               # File dokumentasi Swagger
├── handlers/
│   └── handler.go      # Mengelola request & response HTTP
├── helper/
│   ├── auth.go         # Logika pembuatan token JWT
│   ├── pagination.go   # Logika untuk paginasi
│   └── response.go     # Formatter response JSON standar
├── middleware/
│   └── auth.go         # Middleware untuk validasi token JWT
├── models/             # Definisi struct (request, response, entitas DB)
│   ├── ballance.go     # Model balance
│   ├── category.go     # Model kategori
│   ├── request.go      # Request models (SignUp, Login, Create, Update, etc.)
│   ├── response.go     # Response models (TransactionList, Balance, etc.)
│   ├── transaction.go  # Model transaksi
│   └── user.go         # Model user dengan role
├── repository/
│   ├── repository.go      # Interface untuk interaksi DB
│   └── repository_impl.go # Implementasi interaksi DB
├── services/
│   ├── service.go         # Interface untuk logika bisnis
│   └── service_impl.go    # Implementasi logika bisnis
├── .env                # (Contoh) File variabel lingkungan
├── .gitignore
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── main.go             # Entry point aplikasi dan registrasi rute
└── postman_collection.json # Koleksi Postman untuk testing API
```

## Panduan Menjalankan Proyek

### 1. Persyaratan

-   [Go](https://golang.org/dl/) (versi 1.25 atau lebih baru)
-   [Docker](https://www.docker.com/get-started) & [Docker Compose](https://docs.docker.com/compose/install/)

### 2. Konfigurasi Lingkungan

1.  Buat file `.env` di direktori utama proyek.
2.  Salin konten dari contoh di bawah dan sesuaikan dengan konfigurasi Anda.

```env
# Konfigurasi Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_crud_api
DB_SSLMODE=disable

# Konfigurasi Aplikasi
API_PORT=8080
SECRET_KEY=your_jwt_secret_key_here # Ganti dengan secret key yang kuat untuk JWT
```

### 3. Menjalankan dengan Docker (Direkomendasikan)

Cara termudah untuk menjalankan proyek ini adalah dengan menggunakan Docker Compose. Perintah ini akan membangun image untuk aplikasi Go dan menjalankan container untuk aplikasi serta database PostgreSQL.

```bash
docker-compose up --build
```

Aplikasi akan berjalan dan dapat diakses di `http://localhost:8080`.

### 4. Menjalankan Secara Lokal (Tanpa Docker)

1.  **Jalankan Database**: Pastikan Anda memiliki instance PostgreSQL yang sedang berjalan dan konfigurasikan file `.env` untuk terhubung ke sana. Anda bisa menjalankan PostgreSQL menggunakan Docker secara terpisah:
    ```bash
    docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=go_crud_api -p 5432:5432 -d postgres
    ```

2.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Jalankan Aplikasi**:
    ```bash
    go run main.go
    ```

## Database Migration

Jika Anda mengupgrade dari versi sebelumnya, jalankan migration berikut untuk menambahkan kolom `role` pada tabel users:

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
```

GORM akan otomatis membuat tabel jika belum ada saat aplikasi pertama kali dijalankan.

## Daftar Endpoint API

Semua endpoint berada di bawah prefix `/api/v1`.

### Authentication & Users

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi | Role       |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- | :--------- |
| `POST`   | `/users`                 | Mendaftarkan pengguna baru.                          | Tidak                  | Public     |
| `POST`   | `/login`                 | Login untuk mendapatkan token JWT.                   | Tidak                  | Public     |
| `GET`    | `/users`                 | Mendapatkan detail pengguna yang sedang login.       | Ya                     | All Users  |

### Categories

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi | Role       |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- | :--------- |
| `GET`    | `/categories`            | Mendapatkan daftar kategori (mendukung `limit`, `page`, `q`). | Ya           | All Users  |
| `GET`    | `/categories/:id`        | Mendapatkan detail kategori berdasarkan ID.          | Ya                     | All Users  |
| `POST`   | `/categories`            | Membuat kategori baru.                               | Ya                     | Admin Only |
| `PUT`    | `/categories/:id`        | Memperbarui kategori berdasarkan ID.                 | Ya                     | Admin Only |
| `DELETE` | `/categories/:id`        | Menghapus kategori berdasarkan ID.                   | Ya                     | Admin Only |

### Transactions

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi | Role       |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- | :--------- |
| `POST`   | `/transactions`          | Membuat transaksi baru.                              | Ya                     | All Users  |
| `GET`    | `/transactions`          | Mendapatkan daftar transaksi (mendukung filter `limit`, `page`, `type`, `category_id`, `start_date`, `end_date`, `user_id`*). | Ya | All Users |
| `GET`    | `/transactions/:id`      | Mendapatkan detail transaksi berdasarkan ID.         | Ya                     | All Users  |
| `PUT`    | `/transactions/:id`      | Memperbarui transaksi berdasarkan ID.                | Ya                     | All Users  |
| `DELETE` | `/transactions/:id`      | Menghapus transaksi berdasarkan ID.                  | Ya                     | All Users  |

**Catatan**: 
- User biasa hanya bisa melihat dan mengelola transaksi milik sendiri.
- Admin dapat melihat semua transaksi dari semua user dengan filter `user_id`.
- Default date range: 27 bulan lalu hingga 26 bulan ini.

### Balance

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi | Role       |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- | :--------- |
| `GET`    | `/balance`               | Mendapatkan balance/saldo (mendukung `start_date`, `end_date`). | Ya        | All Users  |

Response menampilkan `total_income`, `total_expense`, dan `balance` (income - expense).

### Admin - User Management

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi | Role       |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- | :--------- |
| `GET`    | `/admin/users`           | Mendapatkan daftar semua user (mendukung `limit`, `page`). | Ya              | Admin Only |
| `POST`   | `/admin/users`           | Membuat user baru.                                   | Ya                     | Admin Only |
| `PUT`    | `/admin/users/:id`       | Memperbarui user berdasarkan ID.                     | Ya                     | Admin Only |
| `DELETE` | `/admin/users/:id`       | Menghapus user berdasarkan ID.                       | Ya                     | Admin Only |

### Contoh Penggunaan Filter Transaksi

```
GET /api/v1/transactions?limit=10&page=1&type=expense&category_id=1&start_date=2026-01-01&end_date=2026-01-31
```

-   `limit=10`: Menampilkan 10 item per halaman.
-   `page=1`: Menampilkan data dari halaman pertama.
-   `type=expense`: Filter transaksi tipe expense (atau `income`).
-   `category_id=1`: Filter berdasarkan kategori dengan ID 1.
-   `start_date=2026-01-01`: Tanggal mulai filter.
-   `end_date=2026-01-31`: Tanggal akhir filter.
-   `user_id=2` (Admin only): Filter transaksi berdasarkan user tertentu.

## Role & Permissions

Sistem ini menggunakan 2 role:

### 👤 User (Default)
- ✅ View categories (read-only)
- ✅ CRUD transaksi sendiri
- ✅ View balance sendiri

### 👨‍💼 Admin
- ✅ CRUD categories (create, update, delete)
- ✅ View semua transaksi dari semua user
- ✅ CRUD user management

### Membuat Admin User

**Saat registrasi**, tambahkan field `role: "admin"`:
```json
POST /api/v1/users
{
  "name": "Admin User",
  "username": "admin",
  "password": "password123",
  "role": "admin"
}
```

**Atau melalui admin panel** (jika sudah ada admin):
```json
POST /api/v1/admin/users
{
  "name": "New Admin",
  "username": "newadmin",
  "password": "password123",
  "role": "admin"
}
```

Default role jika tidak diisi adalah `"user"`.

## Testing dengan Postman

Impor file `postman_collection.json` ke dalam Postman untuk menguji semua endpoint yang tersedia dengan mudah.
- Variabel `{{base_url}}` secara default adalah `http://localhost:8080`.
- Token otentikasi akan otomatis disimpan sebagai variabel koleksi setelah berhasil login.
