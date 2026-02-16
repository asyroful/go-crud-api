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
	}

	router.Run()
	return
}
