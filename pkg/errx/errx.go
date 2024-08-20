// Package errx provides a structured way to create and manage application-specific errors.
// The package is used for creating errors that can be easily converted to HTTP responses or gRPC status codes
// Using this package ensures uniform error handling across different layers of the application.
//
// The errx package includes sub-package "errto":
//
//	errto:
//	   This package provides two main functions (errto.GRPC, errto.HTTP) for converting
//	   ErrorX instances into gRPC status errors and HTTP responses.
//	   These functions are intended to be used in gRPC server handlers and HTTP handlers.
package errx

import (
	"errors"
)

// Compile-time check to ensure ErrorX implements the error interface.
var _ error = (*ErrorX)(nil)

// Type defines the different categories of errors that can be represented by an ErrorX.
type Type int8

const (
	Internal       Type = iota // Internal errors indicate unexpected issues within the application.
	Authentication             // Authentication errors represent issues with user authentication.
	Forbidden                  // Forbidden errors indicate that the user is not allowed to perform the action.
	Validation                 // Validation errors occur when user input does not meet expected criteria.
	NotFound                   // NotFound errors are returned when a requested resource cannot be located.
	Conflict                   // Conflict errors occur when a resource already exists.
)

// New creates a new ErrorX with the given type, message, and code.
// ErrorX instances are used to provide detailed error information that can be easily
// converted into structured error responses for APIs.
func New(errType Type, msg string, code string) *ErrorX {
	return &ErrorX{
		Message: msg,
		Code:    code,
		Type:    errType,
		origin:  errors.New(msg),
	}
}

// ErrorX represents a structured error that can be used within an application.
// It is serializable to JSON, making it suitable for use in API responses and
// can also be converted to gRPC status codes.
type ErrorX struct { // nolint

	// Message provides a human-readable description of the error.
	Message string `json:"message"`

	// Code is a unique identifier for the error, designed for programmatic consumption.
	Code string `json:"code"`

	// Type indicates the category of the error (e.g., Internal, Validation, NotFound).
	Type Type `json:"-"`

	// Details is an optional map of additional details about the error,
	// which can be useful for debugging or providing more context in API responses.
	Details map[string]string `json:"details,omitempty"`

	// origin is used to support error comparison using errors.Is.
	origin error

	// trace is used to store the stack trace of the error.
	trace string
}

// Error implements the error interface for ErrorX.
func (e ErrorX) Error() string {
	return e.Message
}

// Trace returns the stack trace of the error.
func (e *ErrorX) Trace() string {
	return e.trace
}

// WithDetail returns a copy of the ErrorX with an added detail.
// This method is useful for enriching errors with additional context information
// that can help in debugging or providing more informative API responses.
//
// Example usage:
//
//		var NotFoundErr = errx.New(errx.NotFound, "resource not found", errx.CodeNotFound)
//
//		err := repo.FindByID(123)
//		if err != nil {
//			return NotFoundErr.
//				WithDetail("user_id", "123").
//				WithDetail("table", "users").
//	            ... add arbitrary details here ...
//		}
//
// Than in higher layers of application this error can be checked with errors.Is function like this:
//
//	if errors.Is(err, NotFoundErr) {
//		// handle not found error
//	}
func (e *ErrorX) WithDetail(key string, value string) *ErrorX {
	newErr := *e
	if newErr.Details == nil {
		newErr.Details = make(map[string]string)
	}
	newErr.Details[key] = value
	return &newErr
}

// WithCode returns a copy of the ErrorX with the given code.
func (e *ErrorX) WithCode(code string) *ErrorX {
	newErr := *e
	newErr.Code = code
	return &newErr
}

// Is implements the errors.Is interface for ErrorX.
// It allows comparison of two ErrorX instances, returning true if they share the same origin.
func (e *ErrorX) Is(target error) bool {
	t, ok := target.(*ErrorX)
	if !ok {
		return false
	}
	return e.origin == t.origin
}

// GetCode returns the error code of an ErrorX.
// If the error is not an ErrorX, it is considered that
// the error is not properly handled in lower layers and returns Internal code.
func GetCode(err error) string {
	if e, ok := err.(*ErrorX); ok {
		return e.Code
	}

	return CodeInternal
}

// GetType returns the error type of an ErrorX.
// If the error is not an ErrorX, it is considered that
// the error is not properly handled in lower layers and returns Internal type.
func GetType(err error) Type {
	if e, ok := err.(*ErrorX); ok {
		return e.Type
	}

	return Internal
}
