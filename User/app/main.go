package main

import (
	"context"
	"fmt"
	"os"
	"time"
	appservice "userService/application/service"
	"userService/application/usecase"
	_ "userService/docs"
	"userService/infrastructure/http"
	"userService/infrastructure/http/config"
	"userService/infrastructure/http/gin/handler"
	"userService/infrastructure/http/gin/router"
	"userService/infrastructure/persistence/mongodb"
	mongorepository "userService/infrastructure/persistence/mongodb/repository"
	infrastructureservice "userService/infrastructure/service"

	pb "idmService/proto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title User Service API
// @version 1.0
// @description User management service for the POS system. Handles user CRUD operations and ticket purchases.
// @host localhost:12346
// @BasePath /api/user-manager
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	db := mongodb.InitDB()

	userRepo := mongorepository.NewMongoUserRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := userRepo.CreateIndexes(ctx); err != nil {
		fmt.Printf("Warning: Failed to create indexes: %v\n", err)
	}

	idmHost := os.Getenv("IDM_HOST")
	if idmHost == "" {
		idmHost = "localhost"
	}
	idmPort := os.Getenv("IDM_PORT")
	if idmPort == "" {
		idmPort = "50051"
	}

	idmAddress := fmt.Sprintf("%s:%s", idmHost, idmPort)
	idmConn, err := grpc.NewClient(idmAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Warning: Failed to connect to IDM for service auth: %v\n", err)
	}
	defer func() {
		if idmConn != nil {
			idmConn.Close()
		}
	}()

	var serviceAuthClient *infrastructureservice.ServiceAuthClient
	serviceEmail := os.Getenv("SERVICE_EMAIL")
	servicePassword := os.Getenv("SERVICE_PASSWORD")
	if serviceEmail != "" && servicePassword != "" && idmConn != nil {
		idmClient := pb.NewIdentityServiceClient(idmConn)
		serviceAuthClient = infrastructureservice.NewServiceAuthClient(idmClient, serviceEmail, servicePassword)
		fmt.Printf("Service authentication configured for: %s\n", serviceEmail)
	} else {
		fmt.Println("Warning: SERVICE_EMAIL or SERVICE_PASSWORD not set, service-to-service auth disabled")
	}

	eventManagerClient := http.NewEventManagerClient(serviceAuthClient)
	eventManagerService := infrastructureservice.NewEventManagerHTTPAdapter(eventManagerClient)

	authenService, err := infrastructureservice.NewRealAuthenticationService(idmHost, idmPort)
	if err != nil {
		fmt.Printf("Failed to initialize authentication service: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closer, ok := authenService.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	authzService := infrastructureservice.NewDummyAuthorizationService(userRepo)

	// Create UserService layer
	userService := appservice.NewUserService(userRepo)

	userUsecase := usecase.NewUserUsecase(userService, eventManagerService, authenService, authzService)

	serviceURLs := config.NewServiceURLs()

	userHandler := handler.NewGinUserHandler(userUsecase, serviceURLs)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userAPI := r.Group("/api/user-manager")
	router.RegisterUserRoutes(userAPI, userHandler)

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "12346"
	}

	if err := r.Run(":" + port); err != nil {
		fmt.Println(err.Error())
	}
}
