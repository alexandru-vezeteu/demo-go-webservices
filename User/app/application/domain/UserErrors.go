package domain

import "fmt"

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Reason)
}

// NotFoundError represents a user not found error
type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.ID)
}

// AlreadyExistsError represents a duplicate user error
type AlreadyExistsError struct {
	Name string // The field name or value that already exists (e.g., email)
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("user with %s already exists", e.Name)
}

// InternalError represents an internal server error
type InternalError struct {
	Msg string // High-level message
	Err error  // Underlying error
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

// InvalidRequestError represents a malformed request error
type InvalidRequestError struct {
	Reason string
}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Reason)
}

// UnauthorizedError represents an authentication error
type UnauthorizedError struct {
	Reason string
}

func (e *UnauthorizedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("unauthorized: %s", e.Reason)
	}
	return "unauthorized"
}

// ForbiddenError represents an authorization error
type ForbiddenError struct {
	Reason string
}

func (e *ForbiddenError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("forbidden: %s", e.Reason)
	}
	return "forbidden: you don't have permission to access this resource"
}

// DatabaseError represents a database operation error
type DatabaseError struct {
	Operation string // e.g., "create", "update", "delete", "query"
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
