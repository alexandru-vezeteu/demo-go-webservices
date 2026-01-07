package main

import (
	"eventManager/application/service"
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/handler"
	"eventManager/infrastructure/http/gin/router"
	"eventManager/infrastructure/persistence/postgres"
	gormrepository "eventManager/infrastructure/persistence/postgres/gormRepository"
	infrastructureservice "eventManager/infrastructure/service"
	"fmt"
	"os"

	_ "eventManager/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EventManager Service API
// @version 1.0
// @description Event management service for the POS system. Handles events, event packets, tickets, and inclusions.
// @host localhost:12345
// @BasePath /api/event-manager
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	db := postgres.InitDB()

	eventRepo := &gormrepository.GormEventRepository{DB: db}
	eventPacketRepo := &gormrepository.GormEventPacketRepository{DB: db}
	eventPacketInclusionRepo := &gormrepository.GormEventPacketInclusionRepository{DB: db}
	ticketRepo := &gormrepository.GormTicketRepository{DB: db}

	eventService := service.NewEventService(eventRepo, eventPacketInclusionRepo)
	eventPacketService := service.NewEventPacketService(eventPacketRepo, eventPacketInclusionRepo)
	ticketService := service.NewTicketService(ticketRepo, eventRepo, eventPacketRepo, eventPacketInclusionRepo)

	idmHost := os.Getenv("IDM_HOST")
	if idmHost == "" {
		idmHost = "localhost"
	}
	idmPort := os.Getenv("IDM_PORT")
	if idmPort == "" {
		idmPort = "50051"
	}

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

	authzService := infrastructureservice.NewDummyAuthorizationService()

	eventUseCase := usecase.NewEventUseCase(eventRepo, eventService, authenService, authzService)
	eventPacketUseCase := usecase.NewEventPacketUseCase(eventPacketRepo, eventPacketService, authenService, authzService)
	eventPacketInclusionUseCase := usecase.NewEventPacketInclusionUseCase(eventPacketInclusionRepo, eventRepo, eventPacketRepo, authenService, authzService)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo, ticketService, authenService, authzService)

	serviceURLs := config.NewServiceURLs()

	eventHandler := handler.NewGinEventHandler(eventUseCase, serviceURLs)
	eventPacketHandler := handler.NewGinEventPacketHandler(eventPacketUseCase, serviceURLs)
	eventPacketInclusionHandler := handler.NewGinEventPacketInclusionHandler(eventPacketInclusionUseCase, serviceURLs)
	ticketHandler := handler.NewGinTicketHandler(ticketUseCase, serviceURLs)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	eventAPI := r.Group("/api/event-manager")
	router.RegisterEventRoutes(eventAPI, eventHandler)
	router.RegisterEventPacketRoutes(eventAPI, eventPacketHandler)
	router.RegisterEventPacketInclusionRoutes(eventAPI, eventPacketInclusionHandler)
	router.RegisterTicketRoutes(eventAPI, ticketHandler)

	port := os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		fmt.Println(err.Error())
	}
}
