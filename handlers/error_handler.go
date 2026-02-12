package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// FormatValidationError memformat error validasi dari validator menjadi pesan yang lebih mudah dibaca.
// Fungsi ini akan mengembalikan pesan error untuk field pertama yang gagal validasi.
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			// Mengembalikan pesan berdasarkan field yang error
			switch fe.Field() {
			case "Title":
				return "Judul tidak boleh kosong"
			case "Author":
				return "Author tidak boleh kosong"
			}
		}
	}
	return "Input tidak valid"
}
