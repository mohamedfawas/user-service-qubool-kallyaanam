package repository

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is returned when a requested resource doesn't exist
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateKey is returned when a unique constraint is violated
	ErrDuplicateKey = errors.New("duplicate key violation")

	// ErrInvalidOperation is returned when an operation can't be performed
	ErrInvalidOperation = errors.New("invalid operation")

	// ErrTransactionFailed is returned when a transaction operation fails
	ErrTransactionFailed = errors.New("transaction failed")
)

// RepositoryError represents a repository-specific error
type RepositoryError struct {
	Err       error
	Operation string
	Entity    string
	Details   string
}

// Error returns the string representation of the error
func (e *RepositoryError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s failed for %s: %v (%s)", e.Operation, e.Entity, e.Err, e.Details)
	}
	return fmt.Sprintf("%s failed for %s: %v", e.Operation, e.Entity, e.Err)
}

// Unwrap returns the underlying error
func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewError creates a new repository error
func NewError(err error, operation, entity, details string) error {
	return &RepositoryError{
		Err:       err,
		Operation: operation,
		Entity:    entity,
		Details:   details,
	}
}
