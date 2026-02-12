package repository

import (
	"go-crud-api/config"
	"go-crud-api/models"
)

type LoanRepository interface {
	Create(loan models.Loan) models.Loan
	FindActiveLoanByBookID(bookID uint) (models.Loan, error)
}

type loanRepository struct{}

func NewLoanRepository() LoanRepository {
	return &loanRepository{}
}

func (r *loanRepository) Create(loan models.Loan) models.Loan {
	config.DB.Create(&loan)
	return loan
}

func (r *loanRepository) FindActiveLoanByBookID(bookID uint) (models.Loan, error) {
	var loan models.Loan
	err := config.DB.Where("book_id = ? AND is_returned = ?", bookID, false).First(&loan).Error
	return loan, err
}