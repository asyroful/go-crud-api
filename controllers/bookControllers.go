package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-crud-api/handlers"
	"go-crud-api/models"
	"go-crud-api/services"
)

type BookController struct {
	bookService services.BookService
}

func NewBookController(bookService services.BookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

// FindBooks godoc
// @Summary      Ambil semua buku
// @Description  Mengambil daftar semua buku yang tersimpan di database
// @Tags         books
// @Produce      json
// @Success      200  {object}  map[string][]models.Book
// @Router       /books [get]
func (controller *BookController) FindBooks(c *gin.Context) {
	books, err := controller.bookService.FindAllWithLoanDetails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// CreateBook godoc
// @Summary      Tambah buku baru
// @Description  Menyimpan data buku baru ke database
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      models.Book  true  "Data Buku"
// @Success      200   {object}  models.Book
// @Router       /books [post]
func (controller *BookController) CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMsg := handlers.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	createdBook, err := controller.bookService.Create(input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": createdBook})
}

// DeleteBook godoc
// @Summary      Hapus sebuah buku
// @Description  Hapus sebuah buku berdasarkan id
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /books/{id} [delete]
func (controller *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	err := controller.bookService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak ditemukan!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// @Summary      Pinjam buku
// @Tags         loans
// @Accept       json
// @Param        loan  body  models.Loan  true  "Data Peminjaman"
// @Router       /books/borrow [post]
func (controller *BookController) BorrowBook(c *gin.Context) {
	var input struct {
		BookID       uint   `json:"book_id" binding:"required"`
		BorrowerName string `json:"borrower_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := controller.bookService.BorrowBook(input.BookID, input.BorrowerName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": loan})
}