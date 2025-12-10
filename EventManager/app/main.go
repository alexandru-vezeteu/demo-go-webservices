package main

import (
	"eventManager/application/usecase"
	"eventManager/application/service"
	"eventManager/infrastructure/http/gin/handler"
	"eventManager/infrastructure/http/gin/router"
	"eventManager/infrastructure/persistence/postgres"
	gormrepository "eventManager/infrastructure/persistence/postgres/gormRepository"
	"fmt"
	"os"

	_ "eventManager/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Event Service API
// @version         1.0

// @BasePath  /api/event-manager
func main() {
	db := postgres.InitDB()

	// Initialize repositories
	eventRepo := &gormrepository.GormEventRepository{DB: db}
	eventPacketRepo := &gormrepository.GormEventPacketRepository{DB: db}
	eventPacketInclusionRepo := &gormrepository.GormEventPacketInclusionRepository{DB: db}
	ticketRepo := &gormrepository.GormTicketRepository{DB: db}

	// Initialize services with dependencies
	eventService := service.NewEventService(eventRepo, eventPacketInclusionRepo)
	eventPacketService := service.NewEventPacketService(eventPacketRepo, eventPacketInclusionRepo)
	eventPacketInclusionService := service.NewEventPacketInclusionService(eventPacketInclusionRepo, eventRepo, eventPacketRepo)
	ticketService := service.NewTicketService(ticketRepo, eventRepo, eventPacketRepo, eventPacketInclusionRepo)

	// Initialize use cases
	eventUseCase := usecase.NewEventUseCase(eventRepo, eventService)
	eventPacketUseCase := usecase.NewEventPacketUseCase(eventPacketRepo, eventPacketService)
	eventPacketInclusionUseCase := usecase.NewEventPacketInclusionUseCase(eventPacketInclusionRepo, eventPacketInclusionService)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo, ticketService)

	// Initialize handlers
	eventHandler := handler.NewGinEventHandler(eventUseCase)
	eventPacketHandler := handler.NewGinEventPacketHandler(eventPacketUseCase)
	eventPacketInclusionHandler := handler.NewGinEventPacketInclusionHandler(eventPacketInclusionUseCase)
	ticketHandler := handler.NewGinTicketHandler(ticketUseCase)

	r := gin.Default()
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

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
