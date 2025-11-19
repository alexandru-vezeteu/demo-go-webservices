package hateoas

type State string

type Transition struct {
	From    State  `yaml:"from"`
	To      State  `yaml:"to"`
	Rel     string `yaml:"rel"`
	Service string `yaml:"service"`
	Action  string `yaml:"action"`
	Method  string `yaml:"method"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel,omitempty"`
	Method string `json:"method"`
	Title  string `json:"title,omitempty"`
}

type Links map[string]Link

type StateMachineConfig struct {
	Name        string             `yaml:"name"`
	States      []string           `yaml:"states"`
	Transitions []TransitionConfig `yaml:"transitions"`
}

type TransitionConfig struct {
	From    string `yaml:"from"`
	To      string `yaml:"to"`
	Rel     string `yaml:"rel"`
	Service string `yaml:"service"`
	Action  string `yaml:"action"`
	Method  string `yaml:"method"`
	Title   string `yaml:"title"`
}

type ServiceConfig struct {
	Services map[string]ServiceDefinition `yaml:"services"`
}

type ServiceDefinition struct {
	URL    string            `yaml:"url"`
	Routes map[string]string `yaml:"routes"`
}
