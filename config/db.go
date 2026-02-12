// config/db.go
package config

import (
	"fmt"
	"go-crud-api/models"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL")
	// Baris di bawah ini untuk debug, silakan hapus jika sudah berhasil
	fmt.Println("Menghubungkan dengan DSN:", dsn) 

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke database!")
	}

	database.AutoMigrate(&models.Book{}, &models.Loan{})
	DB = database
}