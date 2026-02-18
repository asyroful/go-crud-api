package services

import (
	"go-crud-api/models"
)

type Service interface {
	CreateUser(req models.RequestSignUp) (user models.User, err error)
	GetUserById(req models.RequestGetUserById) (user models.User, err error)
	Login(req models.RequestLogin) (response models.ResponseLogin, err error)
	CreateCategory(req models.RequestCreateCategory) (category models.Category, err error)
	GetCategories(req models.RequestGetCategories) (response models.ResponseCategoryList, err error)
	GetCategoryById(req models.RequestGetCategoryById) (category models.Category, err error)
	UpdateCategory(id int, req models.RequestUpdateCategory) (err error)
	DeleteCategory(id int) (err error)
	CreateTransaction(userId int, req models.RequestCreateTransaction) (response models.TransactionResponse, err error)
	GetTransactions(req models.RequestGetTransactions) (response models.ResponseTransactionList, err error)
	GetTransactionById(req models.RequestGetTransactionById, userId int) (response models.TransactionResponse, err error)
	UpdateTransaction(id int, userId int, req models.RequestUpdateTransaction) (response models.TransactionResponse, err error)
	DeleteTransaction(id int, userId int) (err error)
	GetBalance(req models.RequestGetBalance) (response models.ResponseBalance, err error)
}
