package repository

import (
	"gorm.io/gorm"
	"go-crud-api/models"
)

type repository struct {}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateUser(db *gorm.DB, user models.User) (err error) {
	err = db.Create(&user).Error
	return
}

func (r *repository) FindUserById(db *gorm.DB, id int) (user models.User, err error) {
	err = db.Where("id = ?", id).First(&user).Error
	return
}

func (r *repository) FindUserByUsername(db *gorm.DB, username string) (user models.User, err error) {
	err = db.Where("username = ?", username).First(&user).Error
	return
}