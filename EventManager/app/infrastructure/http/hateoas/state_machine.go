package hateoas

import (
	"fmt"
	"sync"
)

type StateMachine struct {
	name        string
	states      map[State]bool
	transitions map[State][]Transition
	mu          sync.RWMutex
}

func NewStateMachine(name string) *StateMachine {
	return &StateMachine{
		name:        name,
		states:      make(map[State]bool),
		transitions: make(map[State][]Transition),
	}
}

func (sm *StateMachine) AddState(state State) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[state] = true
}

func (sm *StateMachine) AddTransition(t Transition) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.states[t.From] {
		return fmt.Errorf("from state '%s' not registered", t.From)
	}
	if !sm.states[t.To] {
		return fmt.Errorf("to state '%s' not registered", t.To)
	}

	sm.transitions[t.From] = append(sm.transitions[t.From], t)
	return nil
}

func (sm *StateMachine) GetAvailableTransitions(currentState State) []Transition {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	transitions := sm.transitions[currentState]
	result := make([]Transition, len(transitions))
	copy(result, transitions)
	return result
}

func (sm *StateMachine) CanTransition(from, to State) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	transitions := sm.transitions[from]
	for _, t := range transitions {
		if t.To == to {
			return true
		}
	}
	return false
}

func (sm *StateMachine) GetName() string {
	return sm.name
}

func (sm *StateMachine) ValidateState(state State) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.states[state] {
		return fmt.Errorf("state '%s' is not valid for state machine '%s'", state, sm.name)
	}
	return nil
}
