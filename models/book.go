// models/book.go
package models

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" gorm:"index" binding:"required" validate:"required"`
	Author string `json:"author" binding:"required" validate:"required"`
	Loans  []Loan `json:"-" gorm:"foreignKey:BookID"`
}