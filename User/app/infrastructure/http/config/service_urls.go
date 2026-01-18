package config

import (
	"fmt"
	"os"
)

type ServiceURLs struct {
	EventManager string
	UserManager  string
}

func NewServiceURLs() *ServiceURLs {
	return &ServiceURLs{
		EventManager: getEventManagerBaseURL(),
		UserManager:  getUserManagerBaseURL(),
	}
}

func getEventManagerBaseURL() string {
	host := os.Getenv("EVENT_MANAGER_HOST")
	port := os.Getenv("EVENT_MANAGER_PORT")

	return fmt.Sprintf("http://%s:%s/api/event-manager", host, port)
}

func getUserManagerBaseURL() string {
	host := os.Getenv("USER_HOST")
	port := os.Getenv("USER_PORT")

	return fmt.Sprintf("http://%s:%s/api/user-manager", host, port)
}
