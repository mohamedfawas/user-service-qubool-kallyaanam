package postgresutils

import (
	"errors"
	"strings"

	"github.com/jackc/pgconn"
)

const (
	// PostgreSQL error codes
	UniqueViolationCode     = "23505" // Unique constraint violation
	ForeignKeyViolationCode = "23503" // Foreign key violation
	CheckViolationCode      = "23514" // Check constraint violation
)

// IsUniqueConstraintViolation checks if the error is a unique constraint violation
func IsUniqueConstraintViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == UniqueViolationCode &&
			(constraintName == "" || strings.Contains(pgErr.ConstraintName, constraintName))
	}
	return false
}

// IsForeignKeyViolation checks if the error is a foreign key constraint violation
func IsForeignKeyViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == ForeignKeyViolationCode &&
			(constraintName == "" || strings.Contains(pgErr.ConstraintName, constraintName))
	}
	return false
}

// IsCheckViolation checks if the error is a check constraint violation
func IsCheckViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == CheckViolationCode &&
			(constraintName == "" || strings.Contains(pgErr.ConstraintName, constraintName))
	}
	return false
}
