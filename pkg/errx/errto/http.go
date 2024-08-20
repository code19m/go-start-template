package errto

import (
	"encoding/json"
	"fmt"
	"go-start-template/pkg/errx"
	"net/http"
)

// errto.HTTP converts an error to an HTTP response and writes it to the provided http.ResponseWriter.
// This function intended for use in HTTP handlers to convert ErrorX instances to HTTP responses.
// If the error is nil, no response is written.
// If the error is a gRPC status error, it is first converted to an ErrorX using errfrom.GRPC.
func HTTP(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*errx.ErrorX); !ok {
		err = errx.Wrap(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode(err))
	writeBody(w, err)
}

func writeBody(w http.ResponseWriter, err error) {
	if e, ok := err.(*errx.ErrorX); ok {
		errJson, marshalErr := json.Marshal(e)
		if marshalErr == nil {
			_, _ = w.Write(errJson)
			return
		}
		err = fmt.Errorf("marshal error: %w. original error: %w", marshalErr, err)
	}
	_, _ = w.Write(
		[]byte(fmt.Sprintf(`{"message": "%s", "code": "%s"}`, err.Error(), errx.CodeInternal)),
	)
}

func httpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if e, ok := err.(*errx.ErrorX); ok {
		switch e.Type {
		case errx.Authentication:
			return http.StatusUnauthorized
		case errx.Forbidden:
			return http.StatusForbidden
		case errx.Validation:
			return http.StatusBadRequest
		case errx.NotFound:
			return http.StatusNotFound
		case errx.Conflict:
			return http.StatusConflict
		case errx.Internal:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}
