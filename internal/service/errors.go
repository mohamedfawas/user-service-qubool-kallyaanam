package service

import (
	"errors"
	"fmt"
)

// Common service errors
var (
	// ErrValidation is returned when input validation fails
	ErrValidation = errors.New("validation error")

	// ErrNotFound is returned when a requested resource doesn't exist
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicate is returned when a resource already exists
	ErrDuplicate = errors.New("resource already exists")

	// ErrUnauthorized is returned when a user doesn't have permission
	ErrUnauthorized = errors.New("unauthorized action")

	// ErrInternal is returned for unexpected errors
	ErrInternal = errors.New("internal service error")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string
	Message string
}

// ServiceError wraps service-specific errors with context
type ServiceError struct {
	Err       error
	Operation string
	Service   string
	Details   string
	Fields    []ValidationError
}

// Error returns a string representation of the error
func (e *ServiceError) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s failed in %s service: %v (validation errors)",
			e.Operation, e.Service, e.Err)
	}

	if e.Details != "" {
		return fmt.Sprintf("%s failed in %s service: %v (%s)",
			e.Operation, e.Service, e.Err, e.Details)
	}

	return fmt.Sprintf("%s failed in %s service: %v",
		e.Operation, e.Service, e.Err)
}

// Unwrap returns the underlying error
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// GetValidationErrors returns the validation errors if any
func (e *ServiceError) GetValidationErrors() []ValidationError {
	return e.Fields
}

// NewError creates a new service error
func NewError(err error, operation, service, details string) error {
	return &ServiceError{
		Err:       err,
		Operation: operation,
		Service:   service,
		Details:   details,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(operation, service string, fields []ValidationError) error {
	return &ServiceError{
		Err:       ErrValidation,
		Operation: operation,
		Service:   service,
		Fields:    fields,
	}
}
