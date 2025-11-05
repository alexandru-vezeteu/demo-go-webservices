package main

import (
	"eventManager/application/controller"
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

	eventRepo := &gormrepository.GormEventRepository{DB: db}
	eventPacketRepo := &gormrepository.GormEventPacketRepository{DB: db}

	eventService := service.NewEventService(eventRepo)
	eventPacketService := service.NewEventPacketService(eventPacketRepo)

	eventController := controller.NewEventController(eventService)
	eventPacketController := controller.NewEventPacketController(eventPacketService)

	eventHandler := handler.NewGinEventHandler(eventController)
	eventPacketHandler := handler.NewGinEventPacketHandler(eventPacketController)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	eventAPI := r.Group("/api/event-manager")
	router.RegisterEventRoutes(eventAPI, eventHandler)
	router.RegisterEventPacketRoutes(eventAPI, eventPacketHandler)

	port := os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
