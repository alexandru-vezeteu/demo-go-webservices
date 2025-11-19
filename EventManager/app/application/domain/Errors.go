package domain

import "fmt"

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Reason)
	}
	return e.Reason
}

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("ID %d not found", e.ID)
}

type AlreadyExistsError struct {
	Name string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("event '%s' already exists", e.Name)
}

type InternalError struct {
	Msg string
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s (cause: %v)", e.Msg, e.Err)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

type UniqueNameError struct {
	Msg string
}

func (e *UniqueNameError) Error() string {
	return fmt.Sprintf("name %s already exists", e.Msg)
}

type ForeignKeyError struct {
	ID int
}

func (e *ForeignKeyError) Error() string {
	return fmt.Sprintf("id %d does not exist", e.ID)
}

type InvalidRequestError struct {
	Reason string
}

func (e *InvalidRequestError) Error() string {
	return e.Reason
}
