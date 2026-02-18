# Tutorial: Menambahkan API CRUD untuk Entitas `Category`

Dokumen ini akan memandu Anda melalui proses penambahan fungsionalitas *Create, Read, Update, Delete* (CRUD) untuk entitas `Category`, termasuk paginasi dan pencarian.

## 1. Struktur Model

Pastikan Anda memiliki definisi model yang diperlukan.

- **`models/category.go`**: Struct utama `Category`.
- **`models/request.go`**:
    - `RequestCreateCategory`: Untuk membuat kategori baru.
    - `RequestGetCategories`: Untuk query list kategori (dengan paginasi dan search).
    - `RequestGetCategoryById`: Untuk mengambil ID dari parameter URI.
    - `RequestUpdateCategory`: Untuk memperbarui nama kategori.
- **`models/response.go`**:
    - `ResponseCategoryList`: Untuk response daftar kategori yang sudah dipaginasi.

## 2. Lapisan Repository

Repository berinteraksi langsung dengan database.

### `repository/repository.go`

Definisikan semua method yang dibutuhkan dalam `Repository` interface.

```go
type Repository interface {
    // ... (method lain)
    GetCategories(db *gorm.DB, name string, pagination models.QueryPagination) (count int64, categories []models.Category, err error)
    GetCategoryById(db *gorm.DB, id int) (category models.Category, err error)
    UpdateCategory(db *gorm.DB, id int, name string) (err error)
    DeleteCategory(db *gorm.DB, id int) (err error)
}
```

### `repository/repository_impl.go`

Implementasikan method-method tersebut.

- **`GetCategories`**: Method ini menangani paginasi dan pencarian berdasarkan nama.
    - Ia pertama-tama menghitung total data yang cocok dengan kriteria pencarian.
    - Kemudian, ia mengambil data sesuai dengan `limit` dan `offset` yang diberikan.
- **`UpdateCategory`**: Memperbarui record kategori berdasarkan ID.
- **`DeleteCategory`**: Menghapus record kategori berdasarkan ID.

## 3. Lapisan Service

Service berisi logika bisnis aplikasi.

### `services/service.go`

Definisikan method yang akan dipanggil oleh handler di dalam `Service` interface.

```go
type Service interface {
    // ... (method lain)
    GetCategories(req models.RequestGetCategories) (response models.ResponseCategoryList, err error)
    GetCategoryById(req models.RequestGetCategoryById) (category models.Category, err error)
    UpdateCategory(id int, req models.RequestUpdateCategory) (err error)
    DeleteCategory(id int) (err error)
}
```

### `services/service_impl.go`

Implementasikan logika bisnisnya.

- **`GetCategories`**: Memanggil `SetPaginationFromQuery` dari helper untuk menyiapkan data paginasi, lalu memanggil repository untuk mendapatkan data.
- **`UpdateCategory`**: Memanggil repository untuk memperbarui data.
- **`DeleteCategory`**: Memanggil repository untuk menghapus data.

## 4. Lapisan Handler

Handler bertanggung jawab untuk mengelola *request* dan *response* HTTP.

### `handlers/handler.go`

Buat fungsi untuk setiap endpoint.

- **`GetCategories`**: Mengambil parameter `q`, `limit`, dan `page` dari query URL untuk pencarian dan paginasi.
- **`GetCategoryById`**: Mengambil `id` dari URI.
- **`UpdateCategory`**: Mengambil `id` dari URI dan data baru dari body request.
- **`DeleteCategory`**: Mengambil `id` dari URI.

Setiap handler akan memanggil method yang sesuai di service dan memformat response (sukses atau error) menggunakan fungsi helper.

## 5. Pendaftaran Rute

Terakhir, daftarkan semua endpoint baru di `main.go`.

### `main.go`

Tambahkan rute untuk `Category` di dalam grup API `v1`. Semua rute ini diamankan menggunakan middleware otentikasi.

```go
// ...
v1 := router.Group("/api/v1")
{
    // ... (rute lain)
    v1.POST("/categories", auth, handler.CreateCategory)
    v1.GET("/categories", auth, handler.GetCategories)
    v1.GET("/categories/:id", auth, handler.GetCategoryById)
    v1.PUT("/categories/:id", auth, handler.UpdateCategory)
    v1.DELETE("/categories/:id", auth, handler.DeleteCategory)
}
// ...
```

Dengan ini, API CRUD untuk `Category` telah selesai diimplementasikan dengan fungsionalitas paginasi dan pencarian.
