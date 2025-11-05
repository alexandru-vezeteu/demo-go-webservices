package domain

import "fmt"

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
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
	id int
}

func (e *ForeignKeyError) Error() string {
	return fmt.Sprintf("id %d does not exist", e.id)
}
