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

func resolveServiceURL(envURL, hostKey, portKey, apiPath string) (string, error) {
	if url := os.Getenv(envURL); url != "" {
		return url, nil
	}

	host := os.Getenv(hostKey)
	port := os.Getenv(portKey)
	if host == "" || port == "" {
		return "", fmt.Errorf("missing %s/%s configuration", hostKey, portKey)
	}

	return fmt.Sprintf("http://%s:%s/api/%s", host, port, apiPath), nil
}

func getEventManagerBaseURL() string {
	url, err := resolveServiceURL("EVENT_MANAGER_URL", "EVENT_MANAGER_HOST", "EVENT_MANAGER_PORT", "event-manager")
	if err != nil {
		panic(err)
	}
	return url
}

func getUserManagerBaseURL() string {
	url, err := resolveServiceURL("USER_MANAGER_URL", "USER_HOST", "USER_PORT", "user-manager")
	if err != nil {
		panic(err)
	}
	return url
}
