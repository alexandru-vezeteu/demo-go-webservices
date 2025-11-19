package hateoas

import (
	"fmt"
	"strings"
)

type LinkGenerator struct {
	stateMachine *StateMachine
	registry     ServiceRegistry
	selfService  string
}

func NewLinkGenerator(sm *StateMachine, registry ServiceRegistry, selfService string) *LinkGenerator {
	return &LinkGenerator{
		stateMachine: sm,
		registry:     registry,
		selfService:  selfService,
	}
}

func (lg *LinkGenerator) GenerateLinks(resourcePath string, currentState State, params map[string]string) (Links, error) {
	links := Links{}

	selfURL, err := lg.registry.GetServiceURL(lg.selfService)
	if err == nil {
		links["self"] = Link{
			Href:   selfURL + resourcePath,
			Rel:    "self",
			Method: "GET",
		}
	}

	transitions := lg.stateMachine.GetAvailableTransitions(currentState)

	for _, t := range transitions {
		link, err := lg.generateLinkFromTransition(t, params)
		if err != nil {
			continue
		}
		links[t.Rel] = link
	}

	return links, nil
}

func (lg *LinkGenerator) GenerateLinksWithFilter(
	resourcePath string,
	currentState State,
	params map[string]string,
	filterFunc func(Transition) bool,
) (Links, error) {
	links := Links{}

	selfURL, err := lg.registry.GetServiceURL(lg.selfService)
	if err == nil {
		links["self"] = Link{
			Href:   selfURL + resourcePath,
			Rel:    "self",
			Method: "GET",
		}
	}

	transitions := lg.stateMachine.GetAvailableTransitions(currentState)

	for _, t := range transitions {
		if filterFunc != nil && !filterFunc(t) {
			continue
		}

		link, err := lg.generateLinkFromTransition(t, params)
		if err != nil {
			continue
		}
		links[t.Rel] = link
	}

	return links, nil
}

func (lg *LinkGenerator) generateLinkFromTransition(t Transition, params map[string]string) (Link, error) {
	actionURL, err := lg.registry.GetActionURL(t.Service, t.Action)
	if err != nil {
		return Link{}, err
	}

	href := actionURL
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		href = strings.ReplaceAll(href, placeholder, value)
	}

	return Link{
		Href:   href,
		Rel:    t.Rel,
		Method: t.Method,
	}, nil
}

func (lg *LinkGenerator) AddStaticLink(links Links, rel, path, method string) error {
	selfURL, err := lg.registry.GetServiceURL(lg.selfService)
	if err != nil {
		return err
	}

	links[rel] = Link{
		Href:   selfURL + path,
		Rel:    rel,
		Method: method,
	}
	return nil
}
