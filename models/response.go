package models

import "time"

// BookDetailResponse adalah struct untuk menampilkan detail buku beserta status peminjamannya.
type BookDetailResponse struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Author       string     `json:"author"`
	IsBorrowed   bool       `json:"is_borrowed"`
	BorrowerName *string    `json:"borrower_name,omitempty"`
	LoanDate     *time.Time `json:"loan_date,omitempty"`
}
