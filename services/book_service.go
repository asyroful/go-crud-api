package services

import (
	"errors"
	"go-crud-api/models"
	"go-crud-api/repository"
	"strconv"
	"strings"
	"time"
)

type BookService interface {
	FindAll() []models.Book
	FindAllWithLoanDetails() ([]models.BookDetailResponse, error)
	Create(book models.Book) (models.Book, error)
	Delete(id string) error
	FindByTitle(title string) (models.Book, error)
	BorrowBook(bookID uint, borrowerName string) (models.Loan, error)
}

type bookService struct {
	bookRepository repository.BookRepository
	loanRepository repository.LoanRepository
}

func NewBookService(r repository.BookRepository, l repository.LoanRepository) BookService {
	return &bookService{
		bookRepository: r,
		loanRepository: l,
	}
}

func (s *bookService) FindAll() []models.Book {
	return s.bookRepository.FindAll()
}

func (s *bookService) FindAllWithLoanDetails() ([]models.BookDetailResponse, error) {
	books, err := s.bookRepository.FindAllWithLoans()
	if err != nil {
		return nil, err
	}

	var response []models.BookDetailResponse
	for _, book := range books {
		detail := models.BookDetailResponse{
			ID:         book.ID,
			Title:      book.Title,
			Author:     book.Author,
			IsBorrowed: false,
		}

		if len(book.Loans) > 0 {
			// Asumsikan hanya ada satu pinjaman aktif per buku
			loan := book.Loans[0]
			detail.IsBorrowed = true
			detail.BorrowerName = &loan.BorrowerName
			detail.LoanDate = &loan.LoanDate
		}
		response = append(response, detail)
	}

	return response, nil
}

func (s *bookService) Create(book models.Book) (models.Book, error) {
	if strings.TrimSpace(book.Title) == "" {
		return models.Book{}, errors.New("Judul tidak boleh kosong")
	}
	if strings.TrimSpace(book.Author) == "" {
		return models.Book{}, errors.New("Author tidak boleh kosong")
	}

	existingBook, err := s.bookRepository.FindByTitle(book.Title)
	if err == nil && existingBook.ID != 0 {
		return models.Book{}, errors.New("Buku dengan judul '" + book.Title + "' sudah ada")
	}

	return s.bookRepository.Create(book), nil
}

func (s *bookService) Delete(id string) error {
	book, err := s.bookRepository.FindByID(id)
	if err != nil {
		return err
	}
	s.bookRepository.Delete(book)
	return nil
}

func (s *bookService) FindByTitle(title string) (models.Book, error) {
	return s.bookRepository.FindByTitle(title)
}

func (s *bookService) BorrowBook(bookID uint, borrowerName string) (models.Loan, error) {
	// Convert uint to string for FindByID
	bookIDStr := strconv.FormatUint(uint64(bookID), 10)
	_, err := s.bookRepository.FindByID(bookIDStr)
	if err != nil {
		return models.Loan{}, errors.New("buku tidak ditemukan")
	}

	activeLoan, _ := s.loanRepository.FindActiveLoanByBookID(bookID)
	if activeLoan.ID != 0 {
		return models.Loan{}, errors.New("buku ini sedang dipinjam dan belum dikembalikan")
	}

	newLoan := models.Loan{
		BookID:       bookID,
		BorrowerName: borrowerName,
		LoanDate:     time.Now(),
		IsReturned:   false,
	}

	return s.loanRepository.Create(newLoan), nil
}
