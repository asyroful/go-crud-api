package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"go-crud-api/models"
)

func ResponseFormater(statusCode int, message string, data interface{}) models.Response {
	return models.Response{
		Code:    statusCode,
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "success",
		Data:    data,
	})
}
