package repository

import (
	"go-crud-api/models"

	"gorm.io/gorm"
)

type repository struct{}

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

func (r *repository) CreateCategory(db *gorm.DB, category models.Category) (models.Category, error) {
	err := db.Create(&category).Error
	return category, err
}

func (r *repository) GetCategories(db *gorm.DB, name string, pagination models.QueryPagination) (count int64, categories []models.Category, err error) {
	query := db.Model(&models.Category{})

	if name != "" {
		searchQuery := "%" + name + "%"
		query = query.Where("LOWER(name) LIKE ?", searchQuery)
	}

	err = query.Count(&count).Error
	if err != nil {
		return
	}

	err = query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&categories).Error
	if err != nil {
		return
	}

	return
}

func (r *repository) GetCategoryById(db *gorm.DB, id int) (category models.Category, err error) {
	err = db.Where("id = ?", id).First(&category).Error
	return
}

func (r *repository) DeleteCategory(db *gorm.DB, id int) (err error) {
	err = db.Where("id = ?", id).Delete(&models.Category{}).Error
	return
}

func (r *repository) UpdateCategory(db *gorm.DB, id int, name string) (err error) {
	err = db.Model(&models.Category{}).Where("id = ?", id).Update("name", name).Error
	return
}

func (r *repository) CreateTransaction(db *gorm.DB, transaction models.Transaction) (models.Transaction, error) {
	err := db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	// Load relations
	err = db.Preload("User").Preload("Category").First(&transaction, transaction.Id).Error
	return transaction, err
}

func (r *repository) GetTransactions(db *gorm.DB, userId int, categoryId int, transactionType string, pagination models.QueryPagination) (count int64, transactions []models.Transaction, err error) {
	query := db.Model(&models.Transaction{})

	if userId != 0 {
		query = query.Where("user_id = ?", userId)
	}

	if categoryId != 0 {
		query = query.Where("category_id = ?", categoryId)
	}

	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	err = query.Count(&count).Error
	if err != nil {
		return
	}

	err = query.Preload("User").Preload("Category").Order("created_at DESC").Limit(pagination.Limit).Offset(pagination.Offset).Find(&transactions).Error
	if err != nil {
		return
	}

	return
}

func (r *repository) GetTransactionById(db *gorm.DB, id int) (transaction models.Transaction, err error) {
	err = db.Preload("User").Preload("Category").Where("id = ?", id).First(&transaction).Error
	return
}

func (r *repository) UpdateTransaction(db *gorm.DB, id int, transaction models.Transaction) (err error) {
	err = db.Model(&models.Transaction{}).Where("id = ?", id).Updates(transaction).Error
	return
}

func (r *repository) DeleteTransaction(db *gorm.DB, id int) (err error) {
	err = db.Where("id = ?", id).Delete(&models.Transaction{}).Error
	return
}
