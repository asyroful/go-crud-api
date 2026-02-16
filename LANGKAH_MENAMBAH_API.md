# Langkah-langkah Menambahkan API Baru (GET, POST, PUT, DELETE)

Dokumen ini menjelaskan alur kerja untuk menambahkan fungsionalitas API baru ke dalam proyek ini, dengan asumsi Anda sudah memiliki struktur folder `models`, `repository`, `services`, dan `handlers`. Kita akan menggunakan contoh entitas `Product` untuk mempermudah.

## 1. Tentukan Kebutuhan di `models`

Setiap data yang diterima dari *request* atau dikirim sebagai *response* harus didefinisikan strukturnya di dalam direktori `models`.

- **`models/product.go`**: Buat file baru untuk mendefinisikan struct utama dari `Product`.
- **`models/request.go`**: Tambahkan struct untuk data yang diterima dari *request* (misalnya `RequestCreateProduct`, `RequestUpdateProduct`).
- **`models/response.go`**: Tambahkan struct untuk data yang akan dikirim sebagai *response* (misalnya `ResponseProduct`, `ResponseProductList`).

**Contoh:**
```go
// models/product.go
type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Price uint
}

// models/request.go
type RequestCreateProduct struct {
    Name  string `json:"name" binding:"required"`
    Price uint   `json:"price" binding:"required"`
}
```

## 2. Buat Fungsi di `repository`

Lapisan *repository* bertanggung jawab untuk berinteraksi langsung dengan database (query, insert, update, delete).

1.  **Definisikan Interface**: Buka `repository/repository.go` dan tambahkan *method* baru di dalam `Repository` interface.

    ```go
    // repository/repository.go
    type Repository interface {
        // ... method user yang sudah ada
        SaveProduct(product models.Product) (models.Product, error)
        FindProductByID(ID int) (models.Product, error)
        FindAllProducts() ([]models.Product, error)
        UpdateProduct(product models.Product) (models.Product, error)
        DeleteProduct(product models.Product) (models.Product, error)
    }
    ```

2.  **Implementasikan Method**: Buka `repository/repository_impl.go` dan implementasikan *method* yang baru saja Anda definisikan.

    ```go
    // repository/repository_impl.go
    func (r *repository) SaveProduct(product models.Product) (models.Product, error) {
        err := r.db.Create(&product).Error
        return product, err
    }

    // ... implementasi method lainnya
    ```

## 3. Tambahkan Logika Bisnis di `services`

Lapisan *service* berisi logika bisnis. Ia memanggil *method* dari *repository* dan mengolah data sebelum diserahkan ke *handler*.

1.  **Definisikan Interface**: Buka `services/service.go` dan tambahkan *method* baru di dalam `Service` interface.

    ```go
    // services/service.go
    type Service interface {
        // ... method user yang sudah ada
        CreateProduct(req models.RequestCreateProduct) (models.Product, error)
        GetProductByID(ID int) (models.Product, error)
        GetAllProducts() ([]models.Product, error)
        UpdateProduct(ID int, req models.RequestCreateProduct) (models.Product, error)
        DeleteProduct(ID int) (models.Product, error)
    }
    ```

2.  **Implementasikan Method**: Buka `services/service_impl.go` dan implementasikan logikanya. Di sinilah Anda memanggil *repository*.

    ```go
    // services/service_impl.go
    func (s *service) CreateProduct(req models.RequestCreateProduct) (models.Product, error) {
        product := models.Product{
            Name:  req.Name,
            Price: req.Price,
        }
        newProduct, err := s.repository.SaveProduct(product)
        return newProduct, err
    }

    // ... implementasi method lainnya
    ```

## 4. Buat `handler` untuk Menerima Request HTTP

*Handler* adalah jembatan antara *request* HTTP yang masuk dengan logika bisnis di *service*.

1.  **Buat Fungsi Handler**: Buka `handlers/handler.go` dan tambahkan fungsi-fungsi baru untuk setiap *endpoint*.

    ```go
    // handlers/handler.go

    // POST /products
    func (h *handler) CreateProduct(c *gin.Context) {
        var req models.RequestCreateProduct
        if err := c.ShouldBindJSON(&req); err != nil {
            // ... handle error
            return
        }

        product, err := h.service.CreateProduct(req)
        if err != nil {
            // ... handle error
            return
        }
        // ... kirim response sukses
    }

    // GET /products/:id
    func (h *handler) GetProductByID(c *gin.Context) {
        // ... implementasi handler
    }

    // GET /products
    func (h *handler) GetAllProducts(c *gin.Context) {
        // ... implementasi handler
    }

    // PUT /products/:id
    func (h *handler) UpdateProduct(c *gin.Context) {
        // ... implementasi handler
    }

    // DELETE /products/:id
    func (h *handler) DeleteProduct(c *gin.Context) {
        // ... implementasi handler
    }
    ```

## 5. Daftarkan Rute di `main.go`

Langkah terakhir adalah mendaftarkan *endpoint* baru Anda ke *router* Gin agar bisa diakses.

1.  **Tambahkan Rute**: Buka `main.go` dan tambahkan rute baru di dalam grup API Anda.

    ```go
    // main.go
    func main() {
        // ... inisialisasi db, repo, service, handler

        router := gin.Default()
        api := router.Group("/api/v1")

        // ... rute user yang sudah ada

        // Rute untuk Product
        api.POST("/products", handler.CreateProduct)
        api.GET("/products/:id", handler.GetProductByID)
        api.GET("/products", handler.GetAllProducts)
        api.PUT("/products/:id", handler.UpdateProduct)
        api.DELETE("/products/:id", handler.DeleteProduct)

        router.Run()
    }
    ```

Dengan mengikuti langkah-langkah ini secara berurutan, Anda dapat menambahkan fungsionalitas API baru dengan rapi dan terstruktur sesuai dengan arsitektur yang ada.
