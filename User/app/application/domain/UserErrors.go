package domain

import "fmt"


type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Reason)
}


type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.ID)
}


type AlreadyExistsError struct {
	Name string 
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("user with %s already exists", e.Name)
}


type InternalError struct {
	Msg string 
	Err error  
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal error: %s (%v)", e.Msg, e.Err)
	}
	return fmt.Sprintf("internal error: %s", e.Msg)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}


type InvalidRequestError struct {
	Reason string
}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Reason)
}


type UnauthorizedError struct {
	Reason string
}

func (e *UnauthorizedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("unauthorized: %s", e.Reason)
	}
	return "unauthorized"
}


type ForbiddenError struct {
	Reason string
}

func (e *ForbiddenError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("forbidden: %s", e.Reason)
	}
	return "forbidden: you don't have permission to access this resource"
}


type DatabaseError struct {
	Operation string 
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
