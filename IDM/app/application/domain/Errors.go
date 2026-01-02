package domain

import "fmt"

type AuthenticationError struct {
	Reason string
}

func (e *AuthenticationError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("authentication failed: %s", e.Reason)
	}
	return "authentication failed"
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

type TokenError struct {
	Reason      string
	Expired     bool
	Blacklisted bool
	Corrupted   bool
}

func (e *TokenError) Error() string {
	if e.Expired {
		return "token has expired"
	}
	if e.Blacklisted {
		return fmt.Sprintf("token is blacklisted: %s", e.Reason)
	}
	if e.Corrupted {
		return fmt.Sprintf("token is corrupted: %s", e.Reason)
	}
	return fmt.Sprintf("token error: %s", e.Reason)
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for '%s': %s", e.Field, e.Reason)
	}
	return fmt.Sprintf("validation failed: %s", e.Reason)
}

type InternalError struct {
	Operation string
	Err       error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal error during %s: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("internal error during %s", e.Operation)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

type ConfigurationError struct {
	Key    string
	Reason string
}

func (e *ConfigurationError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("configuration error for '%s': %s", e.Key, e.Reason)
	}
	return fmt.Sprintf("configuration error: '%s' is not set", e.Key)
}
