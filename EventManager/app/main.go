package main

import (
	"eventManager/application/service"
	"eventManager/controller"
	"eventManager/infrastructure/http/gin/handler"
	"eventManager/infrastructure/http/gin/router"
	"eventManager/infrastructure/persistence/postgres"
	gormrepository "eventManager/infrastructure/persistence/postgres/gormRepository"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

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
