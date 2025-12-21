package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"userService/application/usecase"
	"userService/infrastructure/http"
	"userService/infrastructure/http/config"
	"userService/infrastructure/http/gin/handler"
	"userService/infrastructure/http/gin/router"
	"userService/infrastructure/persistence/mongodb"
	mongorepository "userService/infrastructure/persistence/mongodb/repository"
	infrastructureservice "userService/infrastructure/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := mongodb.InitDB()

	userRepo := mongorepository.NewMongoUserRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := userRepo.CreateIndexes(ctx); err != nil {
		fmt.Printf("Warning: Failed to create indexes: %v\n", err)
	}

	eventManagerClient := http.NewEventManagerClient()
	eventManagerService := infrastructureservice.NewEventManagerHTTPAdapter(eventManagerClient)

	authenService := infrastructureservice.NewDummyAuthenticationService()
	authzService := infrastructureservice.NewDummyAuthorizationService()

	userUsecase := usecase.NewUserUsecase(userRepo, eventManagerService, authenService, authzService)

	serviceURLs := config.NewServiceURLs()

	userHandler := handler.NewGinUserHandler(userUsecase, serviceURLs)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
