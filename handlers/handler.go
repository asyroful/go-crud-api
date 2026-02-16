package handlers

import (
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/services"
	"net/http"

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
