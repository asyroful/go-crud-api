package models

import "time"

type Loan struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	BookID       uint      `json:"book_id" binding:"required"`
	BorrowerName string    `json:"borrower_name" binding:"required"`
	LoanDate     time.Time `json:"loan_date"`
	ReturnDate   *time.Time `json:"return_date"`
	IsReturned   bool      `json:"is_returned" gorm:"default:false"`
}