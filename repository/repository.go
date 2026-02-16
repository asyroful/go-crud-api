package repository

import (
	"gorm.io/gorm"
	"go-crud-api/models"
)

type Repository interface {
	CreateUser(db *gorm.DB, user models.User) (err error)
	FindUserById(db *gorm.DB, id int) (user models.User, err error)
	FindUserByUsername(db *gorm.DB, username string) (user models.User, err error)
}