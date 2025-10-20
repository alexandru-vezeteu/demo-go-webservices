package main

import (
	"eventManager/Controller"
	"eventManager/Routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	var port = os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		//panic mai bine?
		port = "8080"
	}
	r := gin.Default()

	api := r.Group("/api/event-manager")

	var eventManager = &Controller.EventController{}
	Routes.AddEventRoutes(api, eventManager)

	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
