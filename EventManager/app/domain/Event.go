package domain

import "fmt"

type Event struct {
	ID          int
	OwnerID     int
	Name        string
	Location    string
	Description string
	Seats       int
}

type eventError struct {
	Code    string
	Message string
	Err     error
}

func (e *eventError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *eventError) Unwrap() error {
	return e.Err
}

func NewEventNotFoundError(id int) *eventError {
	return &eventError{
		Code:    "EVENT_NOT_FOUND",
		Message: fmt.Sprintf("event with ID %d not found", id),
		Err:     nil,
	}
}

func NewEventValidationError(message string) *eventError {
	return &eventError{
		Code:    "EVENT_VALIDATION_ERROR",
		Message: message,
		Err:     nil,
	}
}

func NewEventAlreadyExistsError(name string) *eventError {
	return &eventError{
		Code:    "EVENT_ALREADY_EXISTS",
		Message: fmt.Sprintf("event '%s' already exists", name),
		Err:     nil,
	}
}

func NewEventFullError(id int) *eventError {
	return &eventError{
		Code:    "EVENT_FULL",
		Message: fmt.Sprintf("event %d has reached maximum capacity", id),
		Err:     nil,
	}
}

func NewInternalError(message string, cause error) *eventError {
	return &eventError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Err:     cause,
	}
}
