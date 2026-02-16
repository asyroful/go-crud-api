# Go CRUD API (Gin & GORM)

API RESTful sederhana yang dibangun menggunakan Go dengan framework Gin. Proyek ini mengimplementasikan operasi CRUD (Create, Read, Update, Delete) dengan arsitektur bersih (Handlers, Services, Repository) dan menggunakan GORM sebagai ORM untuk berinteraksi dengan database PostgreSQL.

## Fitur

-   **Manajemen Pengguna**: Registrasi dan Login.
-   **Otentikasi JWT**: Endpoint diamankan menggunakan JSON Web Tokens.
-   **CRUD untuk Kategori**:
    -   Membuat, Membaca, Memperbarui, dan Menghapus kategori.
    -   **Paginasi**: Mendukung paginasi (`limit` & `page`) untuk daftar kategori.
    -   **Pencarian**: Mendukung pencarian berdasarkan nama kategori (`q`).
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
API_SECRET=your_jwt_secret_key # Ganti dengan secret key yang kuat
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

## Daftar Endpoint API

Semua endpoint berada di bawah prefix `/api/v1`.

| Method   | Endpoint                 | Deskripsi                                            | Membutuhkan Otentikasi |
| :------- | :----------------------- | :--------------------------------------------------- | :--------------------- |
| `POST`   | `/users`                 | Mendaftarkan pengguna baru.                          | Tidak                  |
| `POST`   | `/login`                 | Login untuk mendapatkan token JWT.                   | Tidak                  |
| `GET`    | `/users`                 | Mendapatkan detail pengguna yang sedang login.       | Ya                     |
| `POST`   | `/categories`            | Membuat kategori baru.                               | Ya                     |
| `GET`    | `/categories`            | Mendapatkan daftar kategori (mendukung `limit`, `page`, `q`). | Ya                     |
| `GET`    | `/categories/:id`        | Mendapatkan detail kategori berdasarkan ID.          | Ya                     |
| `PUT`    | `/categories/:id`        | Memperbarui kategori berdasarkan ID.                 | Ya                     |
| `DELETE` | `/categories/:id`        | Menghapus kategori berdasarkan ID.                   | Ya                     |

### Contoh Penggunaan Paginasi & Pencarian

```
GET /api/v1/categories?limit=5&page=2&q=baju
```

-   `limit=5`: Menampilkan 5 item per halaman.
-   `page=2`: Menampilkan data dari halaman kedua.
-   `q=baju`: Mencari kategori yang namanya mengandung kata "baju".

## Testing dengan Postman

Impor file `postman_collection.json` ke dalam Postman untuk menguji semua endpoint yang tersedia dengan mudah.
- Variabel `{{base_url}}` secara default adalah `http://localhost:8080`.
- Token otentikasi akan otomatis disimpan sebagai variabel koleksi setelah berhasil login.
