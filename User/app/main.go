package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"userService/application/usecase"
	"userService/infrastructure/grpc"
	"userService/infrastructure/http"
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

	// Initialize IDM gRPC client
	idmClient, err2 := grpc.NewIDMClient()
	if err2 != nil {
		log.Fatalf("Failed to create IDM client: %v", err2)
	}
	defer idmClient.Close()

	// Initialize EventManager HTTP client
	eventManagerClient := http.NewEventManagerClient()

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, idmClient, eventManagerClient)

	// Initialize handlers
	userHandler := handler.NewGinUserHandler(userUsecase)

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
