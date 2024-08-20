package errx

const (
	// Default error codes
	CodeInternal       = "INTERNAL"
	CodeAuthentication = "AUTHENTICATION"
	CodeForbidden      = "FORBIDDEN"
	CodeValidation     = "VALIDATION"
	CodeNotFound       = "NOT_FOUND"
	CodeConflict       = "ALREADY_EXISTS"
)

var (
	// Default errors used for common error scenarios.
	ErrInternal       = New(Internal, "Internal server error", CodeInternal)
	ErrAuthentication = New(Authentication, "Unauthenticated", CodeAuthentication)
	ErrForbidden      = New(Forbidden, "Forbidden", CodeForbidden)
	ErrValidation     = New(Validation, "Validation error", CodeValidation)
	ErrNotFound       = New(NotFound, "Resource not found", CodeNotFound)
	ErrConflict       = New(Conflict, "Resource already exists", CodeConflict)
)


// app codes
var (
	CodeOrderStatusInvalid = "invalid_status"
)

func do () {
	
}
