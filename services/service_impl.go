package services

import (
	"errors"
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	Repository repository.Repository
	Db         *gorm.DB
}

func NewService(repository repository.Repository, db *gorm.DB) Service {
	return &service{Repository: repository, Db: db}
}

func (s *service) GetUserById(req models.RequestGetUserById) (user models.User, err error) {
	user, err = s.Repository.FindUserById(s.Db, req.Id)
	return
}

func (s *service) CreateUser(req models.RequestSignUp) (user models.User, err error) {

	if len(req.Name) < 1 && len(req.Username) < 1 && len(req.Password) < 1 {
		err = errors.New("invalid data requested")
		return
	}

	user, err = s.Repository.FindUserByUsername(s.Db, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if user.Id != 0 {
		err = errors.New("username already used")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	user = models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(passwordHash),
	}

	err = s.Repository.CreateUser(s.Db, user)
	if err != nil {
		return
	}

	return
}

func (s *service) Login(req models.RequestLogin) (response models.ResponseLogin, err error) {
	user, err := s.Repository.FindUserByUsername(s.Db, req.Username)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return
	}

	token, err := helper.GenerateToken(user)
	if err != nil {
		return
	}

	response.User = user
	response.Token = token

	return
}

func (s *service) CreateCategory(req models.RequestCreateCategory) (category models.Category, err error) {
	category = models.Category{
		Name: req.Name,
	}
	category, err = s.Repository.CreateCategory(s.Db, category)
	return
}

func (s *service) GetCategories(req models.RequestGetCategories) (response models.ResponseCategoryList, err error) {
	pagination := helper.SetPaginationFromQuery(req.Limit, req.Page)
	count, categories, err := s.Repository.GetCategories(s.Db, req.Name, pagination)
	if err != nil {
		return
	}

	response = models.ResponseCategoryList{
		Count:    count,
		Page:     pagination.Page,
		Limit:    pagination.Limit,
		Data: 		categories,
	}
	return
}

func (s *service) GetCategoryById(req models.RequestGetCategoryById) (category models.Category, err error) {
	category, err = s.Repository.GetCategoryById(s.Db, req.Id)
	return
}

func (s *service) UpdateCategory(id int, req models.RequestUpdateCategory) (err error) {
	err = s.Repository.UpdateCategory(s.Db, id, req.Name)
	return
}

func (s *service) DeleteCategory(id int) (err error) {
	err = s.Repository.DeleteCategory(s.Db, id)
	return
}
