package hapi

import "fmt"

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("parse error for field %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field string
	Value string
	Err   error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" && e.Value != "" {
		return fmt.Sprintf("validation error for field %q with value %q: %v", e.Field, e.Value, e.Err)
	}
	if e.Field != "" {
		return fmt.Sprintf("validation error for field %q: %v", e.Field, e.Err)
	}
	return fmt.Sprintf("validation error: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *ValidationError) Unwrap() error {
	return e.Err
}

