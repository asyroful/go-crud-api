package repository

import (
	"go-crud-api/config"
	"go-crud-api/models"
)

type BookRepository interface {
	FindAll() []models.Book
	Create(book models.Book) models.Book
	Delete(book models.Book)
	FindByID(id string) (models.Book, error)
	FindByTitle(title string) (models.Book, error)
	FindAllWithLoans() ([]models.Book, error)
}

type bookRepository struct{}

func NewBookRepository() BookRepository {
	return &bookRepository{}
}

func (r *bookRepository) FindAll() []models.Book {
	var books []models.Book
	config.DB.Find(&books)
	return books
}

func (r *bookRepository) Create(book models.Book) models.Book {
	config.DB.Create(&book)
	return book
}

func (r *bookRepository) Delete(book models.Book) {
	config.DB.Delete(&book)
}

func (r *bookRepository) FindByID(id string) (models.Book, error) {
	var book models.Book
	if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (r *bookRepository) FindByTitle(title string) (models.Book, error) {
	var book models.Book
	if err := config.DB.Where("title = ?", title).First(&book).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (r *bookRepository) FindAllWithLoans() ([]models.Book, error) {
	var books []models.Book
	// Menggunakan Preload untuk memuat data pinjaman yang aktif (belum dikembalikan)
	// untuk setiap buku.
	err := config.DB.Preload("Loans", "is_returned = ?", false).Find(&books).Error
	return books, err
}
