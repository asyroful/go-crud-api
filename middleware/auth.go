package middleware

import (
	"errors"
	"fmt"
	"go-crud-api/helper"
	"go-crud-api/models"
	"go-crud-api/services"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type AuthMiddleware interface {
	ValidateToken(service services.Service) gin.HandlerFunc
	RequireRole(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{}
}

func loadSecretKey() []byte {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("SECRET_KEY")

	return []byte(secret)
}

func (a authMiddleware) ValidateToken(service services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			errorMessage := gin.H{"errors": errors.New("invalid token").Error()}

			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := tokenValidator(tokenString)
		if err != nil {
			errorMessage := gin.H{"errors": errors.New("invalid token").Error()}

			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			errorMessage := gin.H{"errors": errors.New("invalid token").Error()}

			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		idString := fmt.Sprintf("%v", claim["id"])
		id, _ := strconv.Atoi(idString)

		user, err := service.GetUserById(models.RequestGetUserById{Id: id})
		if err != nil {
			errorMessage := gin.H{"errors": errors.New("invalid token").Error()}

			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if user.Id == 0 {
			errorMessage := gin.H{"errors": errors.New("invalid token").Error()}

			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("current_user", user)
	}
}

func (a authMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("current_user")
		if !exists {
			errorMessage := gin.H{"errors": "unauthorized: user not found"}
			response := helper.ResponseFormater(http.StatusUnauthorized, "error", errorMessage)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		user := currentUser.(models.User)

		// Check if user's role matches any of the required roles
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			errorMessage := gin.H{"errors": "forbidden: insufficient permissions"}
			response := helper.ResponseFormater(http.StatusForbidden, "error", errorMessage)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		c.Next()
	}
}

func tokenValidator(encodedToken string) (*jwt.Token, error) {

	secret := loadSecretKey()

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return secret, nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
