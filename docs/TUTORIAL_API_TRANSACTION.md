# Tutorial: Menambahkan API CRUD untuk Entitas `Transaction`

Dokumen ini akan memandu Anda melalui proses penambahan fungsionalitas *Create, Read, Update, Delete* (CRUD) untuk entitas `Transaction` yang memiliki relasi dengan `User` dan `Category`, termasuk paginasi dan filter.

## 1. Struktur Model

Pastikan Anda memiliki definisi model yang diperlukan.

### `models/transaction.go`

Model `Transaction` memiliki relasi dengan `User` dan `Category`:

```go
package models

import "time"

type Transaction struct {
    Id         int       `json:"id" gorm:"primaryKey"`
    UserId     int       `json:"user_id"`
    User       User      `json:"user" gorm:"foreignKey:UserId"`
    Amount     float64   `json:"amount"`
    Type       string    `json:"type"`
    CategoryId int       `json:"category_id"`
    Category   Category  `json:"category" gorm:"foreignKey:CategoryId"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

### `models/request.go`

Definisikan struct untuk request:

- `RequestCreateTransaction`: Untuk membuat transaksi baru (UserId otomatis dari JWT token).
- `RequestGetTransactions`: Untuk query list transaksi dengan filter dan paginasi.
- `RequestGetTransactionById`: Untuk mengambil ID dari parameter URI.
- `RequestUpdateTransaction`: Untuk memperbarui transaksi.

```go
type RequestCreateTransaction struct {
    Amount     float64 `json:"amount"`
    Type       string  `json:"type"`
    CategoryId int     `json:"category_id"`
}

type RequestGetTransactions struct {
    UserId     int    `json:"user_id"`
    CategoryId string `json:"category_id"`
    Type       string `json:"type"`
    RequestPagination
}

type RequestGetTransactionById struct {
    Id int `json:"id" uri:"id"`
}

type RequestUpdateTransaction struct {
    Amount     float64 `json:"amount"`
    Type       string  `json:"type"`
    CategoryId int     `json:"category_id"`
}
```

### `models/response.go`

Buat response model yang menampilkan relasi `User` dan `Category`:

```go
type ResponseTransactionList struct {
    Data  []TransactionResponse `json:"data"`
    Count int64                 `json:"count"`
    Page  int                   `json:"page"`
    Limit int                   `json:"limit"`
}

type TransactionResponse struct {
    Id        int                    `json:"id"`
    User      UserSimpleResponse     `json:"user"`
    Amount    float64                `json:"amount"`
    Type      string                 `json:"type"`
    Category  CategorySimpleResponse `json:"category"`
    CreatedAt string                 `json:"created_at"`
    UpdatedAt string                 `json:"updated_at"`
}

type UserSimpleResponse struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}

type CategorySimpleResponse struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}
```

## 2. Lapisan Repository

Repository berinteraksi langsung dengan database dan melakukan preload relasi.

### `repository/repository.go`

Definisikan semua method yang dibutuhkan dalam `Repository` interface:

```go
type Repository interface {
    // ... (method lain)
    CreateTransaction(db *gorm.DB, transaction models.Transaction) (models.Transaction, error)
    GetTransactions(db *gorm.DB, userId int, categoryId int, transactionType string, pagination models.QueryPagination) (count int64, transactions []models.Transaction, err error)
    GetTransactionById(db *gorm.DB, id int) (transaction models.Transaction, err error)
    UpdateTransaction(db *gorm.DB, id int, transaction models.Transaction) (err error)
    DeleteTransaction(db *gorm.DB, id int) (err error)
}
```

### `repository/repository_impl.go`

Implementasikan method-method tersebut dengan **Preload** untuk relasi:

```go
func (r *repository) CreateTransaction(db *gorm.DB, transaction models.Transaction) (models.Transaction, error) {
    err := db.Create(&transaction).Error
    if err != nil {
        return transaction, err
    }
    // Load relations after create
    err = db.Preload("User").Preload("Category").First(&transaction, transaction.Id).Error
    return transaction, err
}

func (r *repository) GetTransactions(db *gorm.DB, userId int, categoryId int, transactionType string, pagination models.QueryPagination) (count int64, transactions []models.Transaction, err error) {
    query := db.Model(&models.Transaction{})

    // Filter by user (dari JWT token)
    if userId != 0 {
        query = query.Where("user_id = ?", userId)
    }

    // Filter by category (optional)
    if categoryId != 0 {
        query = query.Where("category_id = ?", categoryId)
    }

    // Filter by type (optional)
    if transactionType != "" {
        query = query.Where("type = ?", transactionType)
    }

    err = query.Count(&count).Error
    if err != nil {
        return
    }

    // Preload User dan Category untuk mendapatkan relasi
    err = query.Preload("User").Preload("Category").Order("created_at DESC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&transactions).Error
    if err != nil {
        return
    }

    return
}

func (r *repository) GetTransactionById(db *gorm.DB, id int) (transaction models.Transaction, err error) {
    // Preload User dan Category
    err = db.Preload("User").Preload("Category").Where("id = ?", id).First(&transaction).Error
    return
}

func (r *repository) UpdateTransaction(db *gorm.DB, id int, transaction models.Transaction) (err error) {
    err = db.Model(&models.Transaction{}).Where("id = ?", id).Updates(transaction).Error
    return
}

func (r *repository) DeleteTransaction(db *gorm.DB, id int) (err error) {
    err = db.Where("id = ?", id).Delete(&models.Transaction{}).Error
    return
}
```

**Catatan Penting**: Penggunaan `Preload("User")` dan `Preload("Category")` sangat penting untuk memuat data relasi dari database.

## 3. Lapisan Service

Service berisi logika bisnis dan transformasi data.

### `services/service.go`

Definisikan method yang akan dipanggil oleh handler:

```go
type Service interface {
    // ... (method lain)
    CreateTransaction(userId int, req models.RequestCreateTransaction) (response models.TransactionResponse, err error)
    GetTransactions(req models.RequestGetTransactions) (response models.ResponseTransactionList, err error)
    GetTransactionById(req models.RequestGetTransactionById) (response models.TransactionResponse, err error)
    UpdateTransaction(id int, req models.RequestUpdateTransaction) (err error)
    DeleteTransaction(id int) (err error)
}
```

### `services/service_impl.go`

Implementasikan logika bisnis dengan transformasi ke response format:

```go
func (s *service) CreateTransaction(userId int, req models.RequestCreateTransaction) (response models.TransactionResponse, err error) {
    // Validasi: Amount tidak boleh 0 atau negatif
    if req.Amount <= 0 {
        err = errors.New("amount must be greater than 0")
        return
    }

    // Validasi: Type tidak boleh kosong atau hanya spasi
    if req.Type == "" {
        err = errors.New("type is required and cannot be empty")
        return
    }

    hasNonSpace := false
    for i := 0; i < len(req.Type); i++ {
        if req.Type[i] != ' ' {
            hasNonSpace = true
            break
        }
    }
    if !hasNonSpace {
        err = errors.New("type cannot contain only spaces")
        return
    }

    // Validasi: CategoryId tidak boleh 0
    if req.CategoryId <= 0 {
        err = errors.New("category_id is required and must be greater than 0")
        return
    }

    // Validasi: Cek apakah category exists
    _, err = s.Repository.GetCategoryById(s.Db, req.CategoryId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            err = errors.New("category not found")
        }
        return
    }

    transaction := models.Transaction{
        UserId:     userId, // Dari JWT token
        Amount:     req.Amount,
        Type:       req.Type,
        CategoryId: req.CategoryId,
    }
    transaction, err = s.Repository.CreateTransaction(s.Db, transaction)
    if err != nil {
        return
    }

    // Transform ke response format
    response = models.TransactionResponse{
        Id: transaction.Id,
        User: models.UserSimpleResponse{
            Id:   transaction.User.Id,
            Name: transaction.User.Name,
        },
        Amount: transaction.Amount,
        Type:   transaction.Type,
        Category: models.CategorySimpleResponse{
            Id:   transaction.Category.Id,
            Name: transaction.Category.Name,
        },
        CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
    return
}

func (s *service) GetTransactions(req models.RequestGetTransactions) (response models.ResponseTransactionList, err error) {
    pagination := helper.SetPaginationFromQuery(req.Limit, req.Page)
    var categoryId int
    if req.CategoryId != "" {
        categoryId, err = strconv.Atoi(req.CategoryId)
        if err != nil {
            return
        }
    }
    count, transactions, err := s.Repository.GetTransactions(s.Db, req.UserId, categoryId, req.Type, pagination)
    if err != nil {
        return
    }

    // Transform ke response format
    transactionResponses := []models.TransactionResponse{}
    for _, transaction := range transactions {
        transactionResponses = append(transactionResponses, models.TransactionResponse{
            Id: transaction.Id,
            User: models.UserSimpleResponse{
                Id:   transaction.User.Id,
                Name: transaction.User.Name,
            },
            Amount: transaction.Amount,
            Type:   transaction.Type,
            Category: models.CategorySimpleResponse{
                Id:   transaction.Category.Id,
                Name: transaction.Category.Name,
            },
            CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
        })
    }

    response = models.ResponseTransactionList{
        Count: count,
        Page:  pagination.Page,
        Limit: pagination.Limit,
        Data:  transactionResponses,
    }
    return
}

func (s *service) GetTransactionById(req models.RequestGetTransactionById, userId int) (response models.TransactionResponse, err error) {
    transaction, err := s.Repository.GetTransactionById(s.Db, req.Id)
    if err != nil {
        return
    }

    // Validate transaction belongs to user
    if transaction.UserId != userId {
        err = errors.New("unauthorized: transaction does not belong to this user")
        return
    }

    // Transform ke response format
    response = models.TransactionResponse{
        Id: transaction.Id,
        User: models.UserSimpleResponse{
            Id:   transaction.User.Id,
            Name: transaction.User.Name,
        },
        Amount: transaction.Amount,
        Type:   transaction.Type,
        Category: models.CategorySimpleResponse{
            Id:   transaction.Category.Id,
            Name: transaction.Category.Name,
        },
        CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
    return
}

func (s *service) UpdateTransaction(id int, userId int, req models.RequestUpdateTransaction) (response models.TransactionResponse, err error) {
    // Validasi: Amount tidak boleh 0 atau negatif
    if req.Amount <= 0 {
        err = errors.New("amount must be greater than 0")
        return
    }

    // Validasi: Type tidak boleh kosong atau hanya spasi
    if req.Type == "" {
        err = errors.New("type is required and cannot be empty")
        return
    }

    hasNonSpace := false
    for i := 0; i < len(req.Type); i++ {
        if req.Type[i] != ' ' {
            hasNonSpace = true
            break
        }
    }
    if !hasNonSpace {
        err = errors.New("type cannot contain only spaces")
        return
    }

    // Validasi: CategoryId tidak boleh 0
    if req.CategoryId <= 0 {
        err = errors.New("category_id is required and must be greater than 0")
        return
    }

    // Validasi: Cek apakah category exists
    _, err = s.Repository.GetCategoryById(s.Db, req.CategoryId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            err = errors.New("category not found")
        }
        return
    }

    // Check if transaction exists and belongs to user
    existingTransaction, err := s.Repository.GetTransactionById(s.Db, id)
    if err != nil {
        return
    }

    if existingTransaction.UserId != userId {
        err = errors.New("unauthorized: transaction does not belong to this user")
        return
    }

    // Update with map to handle all values including zero values
    updateData := map[string]interface{}{
        "amount":      req.Amount,
        "type":        req.Type,
        "category_id": req.CategoryId,
    }

    err = s.Db.Model(&models.Transaction{}).Where("id = ?", id).Updates(updateData).Error
    if err != nil {
        return
    }

    // Get updated transaction with relations
    updatedTransaction, err := s.Repository.GetTransactionById(s.Db, id)
    if err != nil {
        return
    }

    response = models.TransactionResponse{
        Id: updatedTransaction.Id,
        User: models.UserSimpleResponse{
            Id:   updatedTransaction.User.Id,
            Name: updatedTransaction.User.Name,
        },
        Amount: updatedTransaction.Amount,
        Type:   updatedTransaction.Type,
        Category: models.CategorySimpleResponse{
            Id:   updatedTransaction.Category.Id,
            Name: updatedTransaction.Category.Name,
        },
        CreatedAt: updatedTransaction.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: updatedTransaction.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
    return
}

func (s *service) DeleteTransaction(id int, userId int) (err error) {
    // Check if transaction exists and belongs to user
    transaction, err := s.Repository.GetTransactionById(s.Db, id)
    if err != nil {
        return
    }

    if transaction.UserId != userId {
        err = errors.New("unauthorized: transaction does not belong to this user")
        return
    }

    err = s.Repository.DeleteTransaction(s.Db, id)
    return
}
```

## 4. Lapisan Handler

Handler bertanggung jawab untuk mengelola *request* dan *response* HTTP.

### `handlers/handler.go`

Buat fungsi untuk setiap endpoint:

```go
func (h *Handler) CreateTransaction(c *gin.Context) {
    var request models.RequestCreateTransaction

    err := c.ShouldBindJSON(&request)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
        return
    }

    // Ambil userId dari JWT token
    currentUser := c.MustGet("current_user").(models.User)
    userId := currentUser.Id

    transaction, err := h.Service.CreateTransaction(userId, request)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusInternalServerError, response)
        return
    }

    helper.ResponseSuccess(c, transaction)
}

func (h *Handler) GetTransactions(c *gin.Context) {
    var request models.RequestGetTransactions

    // Ambil userId dari JWT token
    currentUser := c.MustGet("current_user").(models.User)
    request.UserId = currentUser.Id
    request.CategoryId = c.Query("category_id")
    request.Type = c.Query("type")
    request.Limit = c.Query("limit")
    request.Page = c.Query("page")

    transactions, err := h.Service.GetTransactions(request)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusInternalServerError, response)
        return
    }
    helper.ResponseSuccess(c, transactions)
}

func (h *Handler) GetTransactionById(c *gin.Context) {
    var request models.RequestGetTransactionById

    err := c.ShouldBindUri(&request)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
        return
    }

    // Get userId from JWT token
    currentUser := c.MustGet("current_user").(models.User)
    userId := currentUser.Id

    transaction, err := h.Service.GetTransactionById(request, userId)
    if err != nil {
        statusCode := http.StatusNotFound
        if err.Error() == "unauthorized: transaction does not belong to this user" {
            statusCode = http.StatusForbidden
        }
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(statusCode, "error", errorMessage)
        c.AbortWithStatusJSON(statusCode, response)
        return
    }

    helper.ResponseSuccess(c, transaction)
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
    var request models.RequestUpdateTransaction
    var id models.RequestGetTransactionById

    err := c.ShouldBindUri(&id)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
        return
    }

    err = c.ShouldBindJSON(&request)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
        return
    }

    // Get userId from JWT token
    currentUser := c.MustGet("current_user").(models.User)
    userId := currentUser.Id

    transaction, err := h.Service.UpdateTransaction(id.Id, userId, request)
    if err != nil {
        statusCode := http.StatusInternalServerError
        if err.Error() == "unauthorized: transaction does not belong to this user" {
            statusCode = http.StatusForbidden
        }
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(statusCode, "error", errorMessage)
        c.AbortWithStatusJSON(statusCode, response)
        return
    }

    helper.ResponseSuccess(c, transaction)
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
    var id models.RequestGetTransactionById

    err := c.ShouldBindUri(&id)
    if err != nil {
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
        c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
        return
    }

    // Get userId from JWT token
    currentUser := c.MustGet("current_user").(models.User)
    userId := currentUser.Id

    err = h.Service.DeleteTransaction(id.Id, userId)
    if err != nil {
        statusCode := http.StatusInternalServerError
        if err.Error() == "unauthorized: transaction does not belong to this user" {
            statusCode = http.StatusForbidden
        }
        errorMessage := gin.H{"errors": err.Error()}
        response := helper.ResponseFormater(statusCode, "error", errorMessage)
        c.AbortWithStatusJSON(statusCode, response)
        return
    }

    helper.ResponseSuccess(c, gin.H{"message": "transaction deleted successfully"})
}
```

**Catatan Penting**: 
- `UserId` otomatis diambil dari JWT token yang sudah di-decode oleh middleware auth.
- User tidak perlu mengirim `user_id` di request body.

## 5. Pendaftaran Rute

Terakhir, daftarkan semua endpoint baru di `main.go`.

### `main.go`

Tambahkan rute untuk `Transaction` di dalam grup API `v1`. Semua rute ini diamankan menggunakan middleware otentikasi:

```go
// ...
v1 := router.Group("/api/v1")
{
    // ... (rute lain)
    
    v1.POST("/transactions", auth, handler.CreateTransaction)
    v1.GET("/transactions", auth, handler.GetTransactions)
    v1.GET("/transactions/:id", auth, handler.GetTransactionById)
    v1.PUT("/transactions/:id", auth, handler.UpdateTransaction)
    v1.DELETE("/transactions/:id", auth, handler.DeleteTransaction)
}
// ...
```

## 6. Migration Database

Pastikan `Transaction` model sudah terdaftar di auto-migration:

### `config/db.go`

```go
database.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{})
```

## Contoh Request & Response

### POST /api/v1/transactions (Create)

**Request Header:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
    "amount": 50000,
    "type": "expense",
    "category_id": 1
}
```

**Response:**
```json
{
    "code": 200,
    "status": "OK",
    "message": "success",
    "data": {
        "id": 1,
        "user": {
            "id": 1,
            "name": "Test User"
        },
        "amount": 50000,
        "type": "expense",
        "category": {
            "id": 1,
            "name": "Food"
        },
        "created_at": "2026-02-16 10:30:00",
        "updated_at": "2026-02-16 10:30:00"
    }
}
```

### GET /api/v1/transactions (List with filters)

**Request Header:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `limit`: Jumlah data per halaman (optional)
- `page`: Nomor halaman (optional)
- `category_id`: Filter by category (optional)
- `type`: Filter by type (expense/income) (optional)

**URL Example:**
```
GET /api/v1/transactions?limit=5&page=1&type=expense&category_id=1
```

**Response:**
```json
{
    "code": 200,
    "status": "OK",
    "message": "success",
    "data": {
        "count": 10,
        "page": 1,
        "limit": 5,
        "data": [
            {
                "id": 1,
                "user": {
                    "id": 1,
                    "name": "Test User"
                },
                "amount": 50000,
                "type": "expense",
                "category": {
                    "id": 1,
                    "name": "Food"
                },
                "created_at": "2026-02-16 10:30:00",
                "updated_at": "2026-02-16 10:30:00"
            }
        ]
    }
}
```

### GET /api/v1/transactions/:id (Get by ID)

**Request Header:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
    "code": 200,
    "status": "OK",
    "message": "success",
    "data": {
        "id": 1,
        "user": {
            "id": 1,
            "name": "Test User"
        },
        "amount": 50000,
        "type": "expense",
        "category": {
            "id": 1,
            "name": "Food"
        },
        "created_at": "2026-02-16 10:30:00",
        "updated_at": "2026-02-16 10:30:00"
    }
}
```

### PUT /api/v1/transactions/:id (Update)

**Request Header:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
    "amount": 60000,
    "type": "expense",
    "category_id": 2
}
```

**Response:**
```json
{
    "code": 200,
    "status": "OK",
    "message": "success",
    "data": {
        "id": 1,
        "user": {
            "id": 1,
            "name": "Test User"
        },
        "amount": 60000,
        "type": "expense",
        "category": {
            "id": 2,
            "name": "Transport"
        },
        "created_at": "2026-02-16 10:30:00",
        "updated_at": "2026-02-16 11:45:00"
    }
}
```

### DELETE /api/v1/transactions/:id (Delete)

**Request Header:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
    "code": 200,
    "status": "OK",
    "message": "success",
    "data": {
        "message": "transaction deleted successfully"
    }
}
```

## Kesimpulan

Dengan mengikuti tutorial ini, Anda telah berhasil membuat API CRUD untuk `Transaction` dengan fitur:

1. ✅ **Relasi Database**: Transaction memiliki relasi dengan User dan Category
2. ✅ **Auto User ID**: UserId otomatis diambil dari JWT token
3. ✅ **Response dengan Relasi**: Response menampilkan user {id, name} dan category {id, name}
4. ✅ **Filter & Pagination**: Mendukung filter berdasarkan category_id dan type
5. ✅ **Security**: Semua endpoint dilindungi dengan JWT authentication
6. ✅ **Authorization**: User hanya bisa melihat, update, dan delete transaksi miliknya sendiri
7. ✅ **Zero Values Handling**: Update menggunakan map sehingga bisa handle semua nilai termasuk zero values
8. ✅ **Validasi Input Ketat**: 
   - Amount > 0 (tidak boleh 0 atau negatif)
   - Type tidak boleh kosong atau hanya spasi
   - CategoryId wajib diisi dan harus exist di database
   - Category name tidak boleh kosong, hanya spasi, atau duplikat

## Fitur Keamanan Tambahan

### 1. Validasi Kepemilikan Transaction
Semua operasi (Get by ID, Update, Delete) memvalidasi bahwa transaction yang diakses benar-benar milik user yang sedang login. Jika tidak, akan mengembalikan error 403 Forbidden.

### 2. Update dengan Map
Menggunakan map untuk update agar semua nilai (termasuk zero values seperti amount=0) tetap bisa di-update. Ini mengatasi limitasi GORM `Updates()` yang men-skip zero values.

### 3. Response Lengkap pada Update
Update transaction tidak hanya mengembalikan message sukses, tetapi mengembalikan data transaction lengkap dengan relasi user dan category.

### 4. Validasi Input Comprehensive
- **Amount**: Harus lebih besar dari 0
- **Type**: Tidak boleh kosong atau hanya berisi spasi
- **CategoryId**: Wajib diisi dan harus exist di database
- **Category Name**: Tidak boleh kosong, hanya spasi, atau duplikat (case-insensitive)

## Error Messages

### Transaction Validation Errors
- `"amount must be greater than 0"` - Amount harus > 0
- `"type is required and cannot be empty"` - Type wajib diisi
- `"type cannot contain only spaces"` - Type tidak boleh hanya spasi
- `"category_id is required"` - CategoryId wajib diisi
- `"category not found"` - CategoryId tidak ditemukan di database
- `"unauthorized: transaction does not belong to this user"` - Transaction bukan milik user

### Category Validation Errors
- `"category name is required and cannot be empty"` - Name wajib diisi
- `"category name cannot contain only spaces"` - Name tidak boleh hanya spasi
- `"category name already exists"` - Name sudah digunakan (case-insensitive)

API Transaction Anda sekarang siap digunakan! 🎉
