package services

import (
	"errors"
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/repository"
	"strconv"
	"time"

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

	// Set default role to "user" if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	user = models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(passwordHash),
		Role:     role,
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
	// Validasi: Name tidak boleh kosong atau hanya spasi
	if len(req.Name) < 1 || req.Name == "" {
		err = errors.New("category name is required and cannot be empty")
		return
	}

	// Trim spasi dan validasi ulang
	trimmedName := req.Name
	for i := 0; i < len(trimmedName); i++ {
		if trimmedName[i] != ' ' {
			break
		}
		if i == len(trimmedName)-1 {
			err = errors.New("category name cannot contain only spaces")
			return
		}
	}

	hasNonSpace := false
	for i := 0; i < len(req.Name); i++ {
		if req.Name[i] != ' ' {
			hasNonSpace = true
			break
		}
	}
	if !hasNonSpace {
		err = errors.New("category name cannot contain only spaces")
		return
	}

	// Validasi: Cek duplikat
	var existingCategory models.Category
	err = s.Db.Where("LOWER(name) = LOWER(?)", req.Name).First(&existingCategory).Error
	if err == nil {
		err = errors.New("category name already exists")
		return
	}
	if err != gorm.ErrRecordNotFound {
		return
	}

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
		Count: count,
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Data:  categories,
	}
	return
}

func (s *service) GetCategoryById(req models.RequestGetCategoryById) (category models.Category, err error) {
	category, err = s.Repository.GetCategoryById(s.Db, req.Id)
	return
}

func (s *service) UpdateCategory(id int, req models.RequestUpdateCategory) (err error) {
	// Validasi: Name tidak boleh kosong atau hanya spasi
	if req.Name == "" {
		err = errors.New("category name is required and cannot be empty")
		return
	}

	hasNonSpace := false
	for i := 0; i < len(req.Name); i++ {
		if req.Name[i] != ' ' {
			hasNonSpace = true
			break
		}
	}
	if !hasNonSpace {
		err = errors.New("category name cannot contain only spaces")
		return
	}

	// Validasi: Cek duplikat (selain category yang sedang di-update)
	var existingCategory models.Category
	err = s.Db.Where("LOWER(name) = LOWER(?) AND id != ?", req.Name, id).First(&existingCategory).Error
	if err == nil {
		err = errors.New("category name already exists")
		return
	}
	if err != gorm.ErrRecordNotFound {
		return
	}

	err = s.Repository.UpdateCategory(s.Db, id, req.Name)
	return
}

func (s *service) DeleteCategory(id int) (err error) {
	err = s.Repository.DeleteCategory(s.Db, id)
	return
}

func (s *service) CreateTransaction(userId int, req models.RequestCreateTransaction) (response models.TransactionResponse, err error) {
	// Validasi: Amount tidak boleh 0 atau negatif
	if req.Amount <= 0 {
		err = errors.New("amount must be greater than 0")
		return
	}

	// Validasi: Type tidak boleh kosong atau hanya spasi
	if req.Type == "" {
		err = errors.New("type is required and cannot be empty")
		return
	}

	hasNonSpace := false
	for i := 0; i < len(req.Type); i++ {
		if req.Type[i] != ' ' {
			hasNonSpace = true
			break
		}
	}
	if !hasNonSpace {
		err = errors.New("type cannot contain only spaces")
		return
	}

	// Validasi: CategoryId wajib diisi
	if req.CategoryId <= 0 {
		err = errors.New("category_id is required")
		return
	}

	// Validasi: Cek apakah category exists
	_, err = s.Repository.GetCategoryById(s.Db, req.CategoryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.New("category not found")
		}
		return
	}

	transaction := models.Transaction{
		UserId:     userId,
		Amount:     req.Amount,
		Type:       req.Type,
		CategoryId: req.CategoryId,
	}
	transaction, err = s.Repository.CreateTransaction(s.Db, transaction)
	if err != nil {
		return
	}

	response = models.TransactionResponse{
		Id: transaction.Id,
		User: models.UserSimpleResponse{
			Id:   transaction.User.Id,
			Name: transaction.User.Name,
		},
		Amount: transaction.Amount,
		Type:   transaction.Type,
		Category: models.CategorySimpleResponse{
			Id:   transaction.Category.Id,
			Name: transaction.Category.Name,
		},
		CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return
}

func (s *service) GetTransactions(req models.RequestGetTransactions) (response models.ResponseTransactionList, err error) {
	pagination := helper.SetPaginationFromQuery(req.Limit, req.Page)
	var categoryId int
	if req.CategoryId != "" {
		categoryId, err = strconv.Atoi(req.CategoryId)
		if err != nil {
			return
		}
	}

	// Set default date range if not provided (27th of previous month to 26th of current month)
	startDate := req.StartDate
	endDate := req.EndDate
	if startDate == "" || endDate == "" {
		now := time.Now()
		// Calculate start date: 27th of previous month
		if now.Day() >= 27 {
			startDate = time.Date(now.Year(), now.Month(), 27, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		} else {
			prevMonth := now.AddDate(0, -1, 0)
			startDate = time.Date(prevMonth.Year(), prevMonth.Month(), 27, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		}
		// Calculate end date: 26th of current month
		if now.Day() >= 27 {
			nextMonth := now.AddDate(0, 1, 0)
			endDate = time.Date(nextMonth.Year(), nextMonth.Month(), 26, 23, 59, 59, 0, now.Location()).Format("2006-01-02")
		} else {
			endDate = time.Date(now.Year(), now.Month(), 26, 23, 59, 59, 0, now.Location()).Format("2006-01-02")
		}
	}

	count, transactions, err := s.Repository.GetTransactions(s.Db, req.UserId, categoryId, req.Type, startDate, endDate, pagination)
	if err != nil {
		return
	}

	// Transform to response format
	transactionResponses := []models.TransactionResponse{}
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, models.TransactionResponse{
			Id: transaction.Id,
			User: models.UserSimpleResponse{
				Id:   transaction.User.Id,
				Name: transaction.User.Name,
			},
			Amount: transaction.Amount,
			Type:   transaction.Type,
			Category: models.CategorySimpleResponse{
				Id:   transaction.Category.Id,
				Name: transaction.Category.Name,
			},
			CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response = models.ResponseTransactionList{
		Count: count,
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Data:  transactionResponses,
	}
	return
}

func (s *service) GetTransactionById(req models.RequestGetTransactionById, userId int) (response models.TransactionResponse, err error) {
	transaction, err := s.Repository.GetTransactionById(s.Db, req.Id)
	if err != nil {
		return
	}

	// Validate transaction belongs to user
	if transaction.UserId != userId {
		err = errors.New("unauthorized: transaction does not belong to this user")
		return
	}

	response = models.TransactionResponse{
		Id: transaction.Id,
		User: models.UserSimpleResponse{
			Id:   transaction.User.Id,
			Name: transaction.User.Name,
		},
		Amount: transaction.Amount,
		Type:   transaction.Type,
		Category: models.CategorySimpleResponse{
			Id:   transaction.Category.Id,
			Name: transaction.Category.Name,
		},
		CreatedAt: transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return
}

func (s *service) UpdateTransaction(id int, userId int, req models.RequestUpdateTransaction) (response models.TransactionResponse, err error) {
	// Validasi: Amount tidak boleh 0 atau negatif
	if req.Amount <= 0 {
		err = errors.New("amount must be greater than 0")
		return
	}

	// Validasi: Type tidak boleh kosong atau hanya spasi
	if req.Type == "" {
		err = errors.New("type is required and cannot be empty")
		return
	}

	hasNonSpace := false
	for i := 0; i < len(req.Type); i++ {
		if req.Type[i] != ' ' {
			hasNonSpace = true
			break
		}
	}
	if !hasNonSpace {
		err = errors.New("type cannot contain only spaces")
		return
	}

	// Validasi: CategoryId wajib diisi
	if req.CategoryId == "" {
		err = errors.New("category_id is required")
		return
	}

	categoryId, err := strconv.Atoi(req.CategoryId)
	if err != nil {
		err = errors.New("invalid category_id format")
		return
	}

	if categoryId <= 0 {
		err = errors.New("category_id must be greater than 0")
		return
	}

	// Validasi: Cek apakah category exists
	_, err = s.Repository.GetCategoryById(s.Db, categoryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.New("category not found")
		}
		return
	}

	// Check if transaction exists and belongs to user
	existingTransaction, err := s.Repository.GetTransactionById(s.Db, id)
	if err != nil {
		return
	}

	if existingTransaction.UserId != userId {
		err = errors.New("unauthorized: transaction does not belong to this user")
		return
	}

	// Update with map to handle all values including zero values
	updateData := map[string]interface{}{
		"amount":      req.Amount,
		"type":        req.Type,
		"category_id": categoryId,
	}

	err = s.Db.Model(&models.Transaction{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return
	}

	// Get updated transaction with relations
	updatedTransaction, err := s.Repository.GetTransactionById(s.Db, id)
	if err != nil {
		return
	}

	response = models.TransactionResponse{
		Id: updatedTransaction.Id,
		User: models.UserSimpleResponse{
			Id:   updatedTransaction.User.Id,
			Name: updatedTransaction.User.Name,
		},
		Amount: updatedTransaction.Amount,
		Type:   updatedTransaction.Type,
		Category: models.CategorySimpleResponse{
			Id:   updatedTransaction.Category.Id,
			Name: updatedTransaction.Category.Name,
		},
		CreatedAt: updatedTransaction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedTransaction.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return
}

func (s *service) DeleteTransaction(id int, userId int) (err error) {
	// Check if transaction exists and belongs to user
	transaction, err := s.Repository.GetTransactionById(s.Db, id)
	if err != nil {
		return
	}

	if transaction.UserId != userId {
		err = errors.New("unauthorized: transaction does not belong to this user")
		return
	}

	err = s.Repository.DeleteTransaction(s.Db, id)
	return
}

func (s *service) GetBalance(req models.RequestGetBalance) (response models.ResponseBalance, err error) {
	// Set default date range if not provided (27th of previous month to 26th of current month)
	startDate := req.StartDate
	endDate := req.EndDate
	if startDate == "" || endDate == "" {
		now := time.Now()
		// Calculate start date: 27th of previous month
		if now.Day() >= 27 {
			startDate = time.Date(now.Year(), now.Month(), 27, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		} else {
			prevMonth := now.AddDate(0, -1, 0)
			startDate = time.Date(prevMonth.Year(), prevMonth.Month(), 27, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		}
		// Calculate end date: 26th of current month
		if now.Day() >= 27 {
			nextMonth := now.AddDate(0, 1, 0)
			endDate = time.Date(nextMonth.Year(), nextMonth.Month(), 26, 23, 59, 59, 0, now.Location()).Format("2006-01-02")
		} else {
			endDate = time.Date(now.Year(), now.Month(), 26, 23, 59, 59, 0, now.Location()).Format("2006-01-02")
		}
	}

	totalIncome, totalExpense, err := s.Repository.GetBalanceByDateRange(s.Db, req.UserId, startDate, endDate)
	if err != nil {
		return
	}

	response = models.ResponseBalance{
		UserId:       req.UserId,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		StartDate:    startDate,
		EndDate:      endDate,
	}

	return
}

// Admin user management methods
func (s *service) GetAllUsers(req models.RequestGetAllUsers) (response models.ResponseUserList, err error) {
	pagination := helper.SetPaginationFromQuery(req.Limit, req.Page)

	count, users, err := s.Repository.GetAllUsers(s.Db, pagination)
	if err != nil {
		return
	}

	userResponses := []models.UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Role:     user.Role,
		})
	}

	response = models.ResponseUserList{
		Count: count,
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Data:  userResponses,
	}

	return
}

func (s *service) AdminCreateUser(req models.RequestCreateUser) (user models.User, err error) {
	if len(req.Name) < 1 || len(req.Username) < 1 || len(req.Password) < 1 {
		err = errors.New("invalid data requested")
		return
	}

	existingUser, err := s.Repository.FindUserByUsername(s.Db, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if existingUser.Id != 0 {
		err = errors.New("username already used")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return
	}

	// Set default role to "user" if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	user = models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(passwordHash),
		Role:     role,
	}

	err = s.Repository.CreateUser(s.Db, user)
	if err != nil {
		return
	}

	return
}

func (s *service) AdminUpdateUser(id int, req models.RequestUpdateUser) (user models.User, err error) {
	// Check if user exists
	user, err = s.Repository.FindUserById(s.Db, id)
	if err != nil {
		return
	}

	if user.Id == 0 {
		err = errors.New("user not found")
		return
	}

	// Prepare update data
	updateData := models.User{}

	if req.Name != "" {
		updateData.Name = req.Name
	}

	if req.Username != "" {
		// Check if username is already taken by another user
		existingUser, _ := s.Repository.FindUserByUsername(s.Db, req.Username)
		if existingUser.Id != 0 && existingUser.Id != id {
			err = errors.New("username already used")
			return
		}
		updateData.Username = req.Username
	}

	if req.Password != "" {
		passwordHash, errHash := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
		if errHash != nil {
			err = errHash
			return
		}
		updateData.Password = string(passwordHash)
	}

	if req.Role != "" {
		updateData.Role = req.Role
	}

	err = s.Repository.UpdateUser(s.Db, id, updateData)
	if err != nil {
		return
	}

	// Get updated user
	user, err = s.Repository.FindUserById(s.Db, id)
	return
}

func (s *service) AdminDeleteUser(id int) (err error) {
	// Check if user exists
	user, err := s.Repository.FindUserById(s.Db, id)
	if err != nil {
		return
	}

	if user.Id == 0 {
		err = errors.New("user not found")
		return
	}

	err = s.Repository.DeleteUser(s.Db, id)
	return
}
