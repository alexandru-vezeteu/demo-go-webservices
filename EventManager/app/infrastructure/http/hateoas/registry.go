package hateoas

import (
	"fmt"
	"strings"
	"sync"
)

type ServiceRegistry interface {
	GetServiceURL(serviceID string) (string, error)
	RegisterService(serviceID, baseURL string)
	RegisterRoute(serviceID, action, path string)
	GetActionURL(serviceID, action string) (string, error)
}

type InMemoryRegistry struct {
	services map[string]string
	routes   map[string]map[string]string
	mu       sync.RWMutex
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		services: make(map[string]string),
		routes:   make(map[string]map[string]string),
	}
}

func (r *InMemoryRegistry) RegisterService(serviceID, baseURL string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[serviceID] = strings.TrimSuffix(baseURL, "/")
}

func (r *InMemoryRegistry) RegisterRoute(serviceID, action, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.routes[serviceID] == nil {
		r.routes[serviceID] = make(map[string]string)
	}
	r.routes[serviceID][action] = path
}

func (r *InMemoryRegistry) GetServiceURL(serviceID string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	baseURL, ok := r.services[serviceID]
	if !ok {
		return "", fmt.Errorf("service '%s' not found in registry", serviceID)
	}
	return baseURL, nil
}

func (r *InMemoryRegistry) GetActionURL(serviceID, action string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	baseURL, ok := r.services[serviceID]
	if !ok {
		return "", fmt.Errorf("service '%s' not found in registry", serviceID)
	}

	if r.routes[serviceID] == nil {
		return "", fmt.Errorf("no routes registered for service '%s'", serviceID)
	}

	path, ok := r.routes[serviceID][action]
	if !ok {
		return "", fmt.Errorf("action '%s' not found for service '%s'", action, serviceID)
	}

	return baseURL + path, nil
}

func (r *InMemoryRegistry) LoadFromConfig(config *ServiceConfig) error {
	for serviceID, def := range config.Services {
		r.RegisterService(serviceID, def.URL)
		for action, path := range def.Routes {
			r.RegisterRoute(serviceID, action, path)
		}
	}
	return nil
}
