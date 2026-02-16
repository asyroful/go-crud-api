package repository

import (
	"go-crud-api/models"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(db *gorm.DB, user models.User) (err error)
	FindUserById(db *gorm.DB, id int) (user models.User, err error)
	FindUserByUsername(db *gorm.DB, username string) (user models.User, err error)
	CreateCategory(db *gorm.DB, category models.Category) (models.Category, error)
	GetCategories(db *gorm.DB, name string, pagination models.QueryPagination) (count int64, categories []models.Category, err error)
	GetCategoryById(db *gorm.DB, id int) (category models.Category, err error)
	UpdateCategory(db *gorm.DB, id int, name string) (err error)
	DeleteCategory(db *gorm.DB, id int) (err error)
}
