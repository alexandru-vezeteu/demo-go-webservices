package domain

import "fmt"

type EventValidationError struct{ Msg string }

func (e *EventValidationError) Error() string { return e.Msg }

func NewEventValidationError(message string) error {
	return &EventValidationError{Msg: message}
}

type EventNotFoundError struct{ ID int }

func (e *EventNotFoundError) Error() string { return fmt.Sprintf("event with ID %d not found", e.ID) }

func NewEventNotFoundError(id int) error {
	return &EventNotFoundError{ID: id}
}

type EventAlreadyExistsError struct{ Name string }

func (e *EventAlreadyExistsError) Error() string {
	return fmt.Sprintf("event '%s' already exists", e.Name)
}

func NewEventAlreadyExistsError(name string) error {
	return &EventAlreadyExistsError{Name: name}
}

type EventFullError struct{ ID int }

func (e *EventFullError) Error() string {
	return fmt.Sprintf("event %d has reached maximum capacity", e.ID)
}

func NewEventFullError(id int) error {
	return &EventFullError{ID: id}
}

type InternalError struct {
	Msg string
	Err error
}

func (e *InternalError) Error() string { return fmt.Sprintf("%s (cause: %v)", e.Msg, e.Err) }
func (e *InternalError) Unwrap() error { return e.Err }

func NewInternalError(message string, cause error) error {
	return &InternalError{Msg: message, Err: cause}
}

type UniqueNameError struct {
	Msg string
}

func (e *UniqueNameError) Error() string { return e.Msg }

func NewUniqueNameError(message string) error {
	return &UniqueNameError{Msg: message}

}

type EventPacketValidationError struct{ Msg string }

func (e *EventPacketValidationError) Error() string { return e.Msg }

func NewEventPacketValidationError(message string) error {
	return &EventPacketValidationError{Msg: message}
}

type EventPacketNotFoundError struct{ ID int }

func (e *EventPacketNotFoundError) Error() string {
	return fmt.Sprintf("event packet with ID %d not found", e.ID)
}

func NewEventPacketNotFoundError(id int) error {
	return &EventPacketNotFoundError{ID: id}
}

type EventPacketAlreadyExistsError struct{ Name string }

func (e *EventPacketAlreadyExistsError) Error() string {
	return fmt.Sprintf("event packet '%s' already exists", e.Name)
}

func NewEventPacketAlreadyExistsError(name string) error {
	return &EventPacketAlreadyExistsError{Name: name}
}
