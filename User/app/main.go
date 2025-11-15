package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"userService/application/controller"
	"userService/application/service"
	"userService/infrastructure/http/gin/handler"
	"userService/infrastructure/http/gin/router"
	"userService/infrastructure/persistence/mongodb"
	mongorepository "userService/infrastructure/persistence/mongodb/repository"

	"github.com/gin-gonic/gin"
)

// @title           User Service API
// @version         1.0

// @BasePath  /api/user-manager
func main() {
	db := mongodb.InitDB()

	// Initialize repositories
	userRepo := mongorepository.NewMongoUserRepository(db)

	// Create indexes for MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := userRepo.CreateIndexes(ctx); err != nil {
		fmt.Printf("Warning: Failed to create indexes: %v\n", err)
	}

	// Initialize services with dependencies
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	userController := controller.NewUserController(userService)

	// Initialize handlers
	userHandler := handler.NewGinUserHandler(userController)

	r := gin.Default()
	userAPI := r.Group("/api/user-manager")
	router.RegisterUserRoutes(userAPI, userHandler)

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "12346"
	}

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
