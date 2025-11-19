package hateoas

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigLoader struct {
	configDir string
}

func NewConfigLoader(configDir string) *ConfigLoader {
	return &ConfigLoader{
		configDir: configDir,
	}
}

func (cl *ConfigLoader) LoadStateMachine(filename string) (*StateMachine, error) {
	path := cl.configDir + "/" + filename
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state machine config %s: %w", path, err)
	}

	var config StateMachineConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse state machine config %s: %w", path, err)
	}

	return BuildStateMachineFromConfig(&config)
}

func (cl *ConfigLoader) LoadServiceConfig(filename string) (*ServiceConfig, error) {
	path := cl.configDir + "/" + filename
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service config %s: %w", path, err)
	}

	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse service config %s: %w", path, err)
	}

	return &config, nil
}

func BuildStateMachineFromConfig(config *StateMachineConfig) (*StateMachine, error) {
	sm := NewStateMachine(config.Name)

	for _, stateName := range config.States {
		sm.AddState(State(stateName))
	}

	for _, tc := range config.Transitions {
		transition := Transition{
			From:    State(tc.From),
			To:      State(tc.To),
			Rel:     tc.Rel,
			Service: tc.Service,
			Action:  tc.Action,
			Method:  tc.Method,
		}

		if err := sm.AddTransition(transition); err != nil {
			return nil, fmt.Errorf("failed to add transition %s->%s: %w", tc.From, tc.To, err)
		}
	}

	return sm, nil
}

func LoadAndBuildStateMachine(configPath string) (*StateMachine, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state machine config: %w", err)
	}

	var config StateMachineConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse state machine config: %w", err)
	}

	return BuildStateMachineFromConfig(&config)
}

func LoadAndBuildServiceRegistry(configPath string) (ServiceRegistry, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service config: %w", err)
	}

	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse service config: %w", err)
	}

	registry := NewInMemoryRegistry()
	if err := registry.LoadFromConfig(&config); err != nil {
		return nil, fmt.Errorf("failed to load service config: %w", err)
	}

	return registry, nil
}
