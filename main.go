package main

import (
	"go-crud-api/config"
	"go-crud-api/controllers"
	_ "go-crud-api/docs"
	"go-crud-api/repository"
	"go-crud-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title        Go CRUD API
// @version      1.0
// @description  Contoh sederhana CRUD API dengan Go, Gin, dan GORM.

// @host      localhost:8080
// @BasePath  /
func main() {
	godotenv.Load()
	config.ConnectDatabase()

	r := gin.Default()

	loanRepository := repository.NewLoanRepository()
	bookRepository := repository.NewBookRepository()
	bookService := services.NewBookService(bookRepository, loanRepository)
	bookController := controllers.NewBookController(bookService)

	// Route Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes API
	v1 := r.Group("/api/v1") // Menggunakan grouping agar rapi
	{
		// Route Buku
		v1.GET("/books", bookController.FindBooks)
		v1.POST("/books", bookController.CreateBook)
		v1.DELETE("/books/:id", bookController.DeleteBook)

		// Route Peminjaman (Fitur Baru)
		v1.POST("/books/borrow", bookController.BorrowBook)
	}

	r.Run()
}