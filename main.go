package main

import (
	"go-crud-api/config"
	_ "go-crud-api/docs"
	"go-crud-api/handlers"
	"go-crud-api/middleware"
	"go-crud-api/repository"
	"go-crud-api/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	config.ConnectDatabase()

	router := gin.Default()
	repo := repository.NewRepository()
	service := services.NewService(repo, config.DB)
	handler := handlers.NewHandler(service)
	mid := middleware.NewAuthMiddleware()

	auth := mid.ValidateToken(service)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"authorization", "content-type"}
	router.Use(cors.New(corsConfig))

	// Routes API
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", handler.CreateUser)
		v1.GET("/users", auth, handler.GetUserById)
		v1.POST("/login", handler.Login)

		v1.POST("/categories", auth, handler.CreateCategory)
		v1.GET("/categories", auth, handler.GetCategories)
		v1.GET("/categories/:id", auth, handler.GetCategoryById)
		v1.PUT("/categories/:id", auth, handler.UpdateCategory)
		v1.DELETE("/categories/:id", auth, handler.DeleteCategory)

		v1.POST("/transactions", auth, handler.CreateTransaction)
		v1.GET("/transactions", auth, handler.GetTransactions)
		v1.GET("/transactions/:id", auth, handler.GetTransactionById)
		v1.PUT("/transactions/:id", auth, handler.UpdateTransaction)
		v1.DELETE("/transactions/:id", auth, handler.DeleteTransaction)

		v1.GET("/balance", auth, handler.GetBalance)
	}

	router.Run()
}
