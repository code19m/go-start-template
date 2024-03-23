// apperr package provides a way to create application-specific errors.
// It is intended to be used to describe application-specific errors.
// Not infrastructure errors like database connection errors or network errors.
package apperr

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// Compile time check for interface implementation
var _ error = (*AppError)(nil)

type Type string

const (
	Authorization Type = "AUTHORIZATION"
	Forbidden     Type = "FORBIDDEN"
	Validation    Type = "VALIDATION"
	NotFound      Type = "NOT_FOUND"
	Conflict      Type = "CONFLICT"
	Internal      Type = "INTERNAL"

	// Default error codes
	DefaultAuthorizationCode = 1000
	DefaultForbiddenCode     = 2000
	DefaultValidationCode    = 3000
	DefaultNotFoundCode      = 4000
	DefaultConflictCode      = 5000
	DefaultInternalCode      = 6000
)

// New creates a new AppError with the given type and message.
func New(errType Type, msg string, code int) *AppError {
	return &AppError{
		Message: msg,
		Code:    code,
		Type:    errType,

		origin: errors.New(msg),
	}
}

// AppError is a struct that represents an error that can be returned by an application.
// It can be serialized to JSON and used to generate HTTP and gRPC status codes.
type AppError struct {

	// origin is necessary for using errors.Is function.
	origin error

	// Message is a human-readable description of the error.
	// It is intended for an end-user audience and should not contain technical information.
	Message string `json:"message"`

	// Code is a unique identifier for the error.
	// It is intended to be consumed programmatically.
	Code int `json:"code"`

	// Type is an identifier for the error type.
	// It can be used to generate handler status codes.
	// For example, a "NOT_FOUND" error type could be used
	// to generate a 404 HTTP code and a 5 gRPC code
	Type Type `json:"type"`

	// Details is an optional map field of additional details about the error
	// that can be used to generate more detailed error messages.
	// This field is especially useful for debugging.
	Details map[string]any `json:"details,omitempty"`
}

// Error method implements the error interface for AppError.
func (e AppError) Error() string {
	return e.Message
}

// WithDetail returns a new copy of AppError with the given detail information.
// Example usage:
//
//	UserNotFoundErr := apperr.New(apperr.NotFound, "User not found", apperr.DefaultNotFoundCode)
//
//	err := repo.FindUserByID(123)
//
//	if err != nil {
//		return UserNotFoundErr.
//			WithDetail("user_id", 123).
//			WithDetail("some_other_info", "***")
//	}
func (e *AppError) WithDetail(key string, value any) *AppError {
	// Create a new copy of the AppError
	newErr := *e
	if newErr.Details == nil {
		newErr.Details = make(map[string]any)
	}
	newErr.Details[key] = value
	return &newErr
}

// Is method implements the errors.Is interface for AppError.
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.origin == t.origin
}

// HTTPCode returns the HTTP status code for the given error.
func HTTPCode(err error) int {
	if e, ok := err.(*AppError); ok {
		switch e.Type {
		case Authorization:
			return http.StatusUnauthorized
		case Forbidden:
			return http.StatusForbidden
		case Validation:
			return http.StatusBadRequest
		case NotFound:
			return http.StatusNotFound
		case Conflict:
			return http.StatusConflict
		case Internal:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

// GRPCCode returns the gRPC status code for the given error.
func GRPCCode(err error) int {
	if e, ok := err.(*AppError); ok {
		switch e.Type {
		case Authorization:
			return int(codes.Unauthenticated)
		case Forbidden:
			return int(codes.PermissionDenied)
		case Validation:
			return int(codes.InvalidArgument)
		case NotFound:
			return int(codes.NotFound)
		case Conflict:
			return int(codes.AlreadyExists)
		case Internal:
			return int(codes.Internal)
		}
	}
	return int(codes.Internal)
}
