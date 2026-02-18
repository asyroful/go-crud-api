package handlers

import (
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) fetchUser(c *gin.Context) {
	currentUser := c.MustGet("current_user").(models.User)
	if currentUser.Id < 1 {
		errorMessage := gin.H{"errors": "not identifier included"}

		response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	helper.ResponseSuccess(c, currentUser)
}

func (h *Handler) CreateUser(c *gin.Context) {

	var request models.RequestSignUp

	err := c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.Service.CreateUser(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	userResponse := models.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
	}

	helper.ResponseSuccess(c, userResponse)
}

func (h *Handler) GetUserById(c *gin.Context) {

	currentUser := c.MustGet("current_user").(models.User)

	userResponse := models.UserResponse{
		Id:       currentUser.Id,
		Name:     currentUser.Name,
		Username: currentUser.Username,
	}

	response := helper.ResponseFormater(http.StatusOK, "success", userResponse)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Login(c *gin.Context) {

	var request models.RequestLogin

	err := c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	loginResult, err := h.Service.Login(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseFormater(http.StatusNotFound, "error", errorMessage)

		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	loginResponse := models.LoginResponse{
		Id:       loginResult.User.Id,
		Name:     loginResult.User.Name,
		Username: loginResult.User.Username,
		Token:    loginResult.Token,
	}

	helper.ResponseSuccess(c, loginResponse)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var request models.RequestCreateCategory

	err := c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	category, err := h.Service.CreateCategory(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, category)
}

func (h *Handler) GetCategories(c *gin.Context) {
	var request models.RequestGetCategories
	request.Name = c.Query("q")
	request.Limit = c.Query("limit")
	request.Page = c.Query("page")

	categories, err := h.Service.GetCategories(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, categories)
}

func (h *Handler) GetCategoryById(c *gin.Context) {
	var request models.RequestGetCategoryById

	err := c.ShouldBindUri(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	category, err := h.Service.GetCategoryById(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusNotFound, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	helper.ResponseSuccess(c, category)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var request models.RequestUpdateCategory
	var id models.RequestGetCategoryById

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.Service.UpdateCategory(id.Id, request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{"message": "category updated successfully"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	var id models.RequestGetCategoryById

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.Service.DeleteCategory(id.Id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{"message": "category deleted successfully"})
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var request models.RequestCreateTransaction

	err := c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("current_user").(models.User)
	userId := currentUser.Id

	transaction, err := h.Service.CreateTransaction(userId, request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, transaction)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	var request models.RequestGetTransactions

	currentUser := c.MustGet("current_user").(models.User)

	// If admin, allow querying all users' transactions
	// If regular user, only show their own transactions
	if currentUser.Role == "admin" {
		// Admin can optionally filter by user_id via query param
		userIdQuery := c.Query("user_id")
		if userIdQuery != "" {
			// Parse and set user_id if provided
			userId, err := strconv.Atoi(userIdQuery)
			if err == nil {
				request.UserId = userId
			}
		}
		// If no user_id provided, UserId will be 0 and repository will return all transactions
	} else {
		// Regular user can only see their own transactions
		request.UserId = currentUser.Id
	}

	request.CategoryId = c.Query("category_id")
	request.Type = c.Query("type")
	request.StartDate = c.Query("start_date")
	request.EndDate = c.Query("end_date")
	request.Limit = c.Query("limit")
	request.Page = c.Query("page")

	transactions, err := h.Service.GetTransactions(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	helper.ResponseSuccess(c, transactions)
}

func (h *Handler) GetTransactionById(c *gin.Context) {
	var request models.RequestGetTransactionById

	err := c.ShouldBindUri(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Get userId from JWT token
	currentUser := c.MustGet("current_user").(models.User)
	userId := currentUser.Id

	transaction, err := h.Service.GetTransactionById(request, userId)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "unauthorized: transaction does not belong to this user" {
			statusCode = http.StatusForbidden
		}
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(statusCode, "error", errorMessage)
		c.AbortWithStatusJSON(statusCode, response)
		return
	}

	helper.ResponseSuccess(c, transaction)
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	var request models.RequestUpdateTransaction
	var id models.RequestGetTransactionById

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Get userId from JWT token
	currentUser := c.MustGet("current_user").(models.User)
	userId := currentUser.Id

	transaction, err := h.Service.UpdateTransaction(id.Id, userId, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "unauthorized: transaction does not belong to this user" {
			statusCode = http.StatusForbidden
		}
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(statusCode, "error", errorMessage)
		c.AbortWithStatusJSON(statusCode, response)
		return
	}

	helper.ResponseSuccess(c, transaction)
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	var id models.RequestGetTransactionById

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Get userId from JWT token
	currentUser := c.MustGet("current_user").(models.User)
	userId := currentUser.Id

	err = h.Service.DeleteTransaction(id.Id, userId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "unauthorized: transaction does not belong to this user" {
			statusCode = http.StatusForbidden
		}
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(statusCode, "error", errorMessage)
		c.AbortWithStatusJSON(statusCode, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{"message": "transaction deleted successfully"})
}

func (h *Handler) GetBalance(c *gin.Context) {
	var request models.RequestGetBalance

	currentUser := c.MustGet("current_user").(models.User)
	request.UserId = currentUser.Id
	request.StartDate = c.Query("start_date")
	request.EndDate = c.Query("end_date")

	balance, err := h.Service.GetBalance(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, balance)
}

// Admin user management handlers
func (h *Handler) GetAllUsers(c *gin.Context) {
	var request models.RequestGetAllUsers
	request.Limit = c.Query("limit")
	request.Page = c.Query("page")

	users, err := h.Service.GetAllUsers(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusInternalServerError, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	helper.ResponseSuccess(c, users)
}

func (h *Handler) AdminCreateUser(c *gin.Context) {
	var request models.RequestCreateUser

	err := c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.Service.AdminCreateUser(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusBadRequest, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{
		"id":       user.Id,
		"name":     user.Name,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *Handler) AdminUpdateUser(c *gin.Context) {
	var id models.RequestDeleteUser

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	var request models.RequestUpdateUser
	err = c.ShouldBindJSON(&request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.Service.AdminUpdateUser(id.Id, request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusBadRequest, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{
		"id":       user.Id,
		"name":     user.Name,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	var id models.RequestDeleteUser

	err := c.ShouldBindUri(&id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusUnprocessableEntity, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.Service.AdminDeleteUser(id.Id)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseFormater(http.StatusBadRequest, "error", errorMessage)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	helper.ResponseSuccess(c, gin.H{"message": "user deleted successfully"})
}
