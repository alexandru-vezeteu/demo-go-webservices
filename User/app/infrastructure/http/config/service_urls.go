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
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		port = "12345"
	}
	return fmt.Sprintf("http://%s:%s/api/event-manager", host, port)
}

func getUserManagerBaseURL() string {
	host := os.Getenv("USER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "12346"
	}
	return fmt.Sprintf("http://%s:%s/api/user-manager", host, port)
}
