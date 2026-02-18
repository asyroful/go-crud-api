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
	adminOnly := mid.RequireRole("admin")

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

		// Category routes - admin can CRUD, users can only read
		v1.GET("/categories", auth, handler.GetCategories)
		v1.GET("/categories/:id", auth, handler.GetCategoryById)
		v1.POST("/categories", auth, adminOnly, handler.CreateCategory)
		v1.PUT("/categories/:id", auth, adminOnly, handler.UpdateCategory)
		v1.DELETE("/categories/:id", auth, adminOnly, handler.DeleteCategory)

		// Transaction routes - users can CRUD their own, admin can see all
		v1.POST("/transactions", auth, handler.CreateTransaction)
		v1.GET("/transactions", auth, handler.GetTransactions)
		v1.GET("/transactions/:id", auth, handler.GetTransactionById)
		v1.PUT("/transactions/:id", auth, handler.UpdateTransaction)
		v1.DELETE("/transactions/:id", auth, handler.DeleteTransaction)

		v1.GET("/balance", auth, handler.GetBalance)

		// Admin user management routes
		v1.GET("/admin/users", auth, adminOnly, handler.GetAllUsers)
		v1.POST("/admin/users", auth, adminOnly, handler.AdminCreateUser)
		v1.PUT("/admin/users/:id", auth, adminOnly, handler.AdminUpdateUser)
		v1.DELETE("/admin/users/:id", auth, adminOnly, handler.AdminDeleteUser)
	}

	router.Run()
}
