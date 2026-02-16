package services

import (
	"errors"
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/repository"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	Repository repository.Repository
	Db				 *gorm.DB
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
