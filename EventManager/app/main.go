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
	eventService := service.NewEventService(eventRepo)
	eventController := controller.NewEventController(eventService)
	eventHandler := handler.NewGinEventHandler(eventController)

	r := gin.Default()
	eventAPI := r.Group("/api/event-manager")
	router.RegisterEventRoutes(eventAPI, eventHandler)

	port := os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
