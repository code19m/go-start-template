package apperr

type Type string

const (
	Authorization          Type = "AUTHORIZATION"
	Forbidden              Type = "FORBIDDEN"
	Validation             Type = "VALIDATION"
	NotFound               Type = "NOT_FOUND"
	Conflict               Type = "CONFLICT"
	PayloadTooLarge        Type = "PAYLOAD_TOO_LARGE"
	UnsupportedPayloadType Type = "UNSUPPORTED_PAYLOAD_TYPE"
	TimeOut                Type = "TIMEOUT"
	Internal               Type = "INTERNAL"
)

func NewAppError(errType Type, msg string) error {
	return &AppError{
		Type:    errType,
		Message: msg,
	}
}

type AppError struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}
