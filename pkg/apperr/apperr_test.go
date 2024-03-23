package apperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	var err error = New(Authorization, "Unauthorized", DefaultAuthorizationCode)

	appErr, ok := err.(*AppError)
	if !ok {
		t.Errorf("Expected AppError, got %T", err)
	}

	if appErr.Message != "Unauthorized" {
		t.Errorf("Expected message 'Unauthorized', got %s", appErr.Message)
	}

	if appErr.Code != DefaultAuthorizationCode {
		t.Errorf("Expected code %d, got %d", DefaultAuthorizationCode, appErr.Code)
	}

	if appErr.Type != Authorization {
		t.Errorf("Expected type Authorization, got %s", appErr.Type)
	}

	if appErr.Error() != "Unauthorized" {
		t.Errorf("Expected error message 'Unauthorized', got %s", appErr.Error())
	}
}

func TestWithDetail(t *testing.T) {
	err := New(Validation, "Validation error", DefaultValidationCode)
	err = err.WithDetail("field", "username")

	if len(err.Details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(err.Details))
	}

	if err.Details["field"] != "username" {
		t.Errorf("Expected detail 'username' for 'field', got %s", err.Details["field"])
	}
}

func TestHTTPCode(t *testing.T) {
	tests := []struct {
		err    error
		status int
	}{
		{New(Authorization, "Unauthorized", DefaultAuthorizationCode), http.StatusUnauthorized},
		{New(Forbidden, "Forbidden", DefaultForbiddenCode), http.StatusForbidden},
		{New(Validation, "Validation error", DefaultValidationCode), http.StatusBadRequest},
		{New(NotFound, "Not found", DefaultNotFoundCode), http.StatusNotFound},
		{New(Conflict, "Conflict", DefaultConflictCode), http.StatusConflict},
		{New(Internal, "Internal error", DefaultInternalCode), http.StatusInternalServerError},
		{errors.New("unknown error"), http.StatusInternalServerError},
	}

	for _, test := range tests {
		if HTTPCode(test.err) != test.status {
			t.Errorf("Expected HTTP status %d for error, got %d", test.status, HTTPCode(test.err))
		}
	}
}

func TestGRPCCode(t *testing.T) {
	tests := []struct {
		err  error
		code int
	}{
		{New(Authorization, "Unauthorized", DefaultAuthorizationCode), int(codes.Unauthenticated)},
		{New(Forbidden, "Forbidden", DefaultForbiddenCode), int(codes.PermissionDenied)},
		{New(Validation, "Validation error", DefaultValidationCode), int(codes.InvalidArgument)},
		{New(NotFound, "Not found", DefaultNotFoundCode), int(codes.NotFound)},
		{New(Conflict, "Conflict", DefaultConflictCode), int(codes.AlreadyExists)},
		{New(Internal, "Internal error", DefaultInternalCode), int(codes.Internal)},
		{errors.New("unknown error"), int(codes.Internal)},
	}

	for _, test := range tests {
		if GRPCCode(test.err) != test.code {
			t.Errorf("Expected gRPC code %d for error, got %d", test.code, GRPCCode(test.err))
		}
	}
}

func TestUnwrapIsAs(t *testing.T) {
	originalErr := New(Validation, "Validation error", DefaultValidationCode)
	wrappedErr := fmt.Errorf("wrapped error: %w", originalErr)

	// Test Unwrap
	if errors.Unwrap(wrappedErr) != originalErr {
		t.Errorf("Expected Unwrap to return the original error")
	}

	// Test Is
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("Expected Is to return true for the original error")
	}
	if errors.Is(wrappedErr, errors.New("Validation error")) {
		t.Errorf("Expected Is to return false for unknown error")
	}

	// Test As
	var target *AppError
	if !errors.As(wrappedErr, &target) {
		t.Errorf("Expected As to return true and set target")
	}

	if target != originalErr {
		t.Errorf("Expected target to be the same as appErr")
	}
}

func TestUnwrapIsAsWithDetails(t *testing.T) {
	originalErr := New(Validation, "Validation error", DefaultValidationCode)
	detailedErr := originalErr.WithDetail("field", "username")

	wrappedErr := fmt.Errorf("wrapped error: %w", detailedErr)

	// Test Unwrap
	unwrappedErr := errors.Unwrap(wrappedErr)
	if unwrappedErr != detailedErr {
		t.Errorf("Expected Unwrap to return the detailed error")
	}

	// Test Is with detailed and original errors
	if !errors.Is(wrappedErr, detailedErr) {
		t.Errorf("Expected Is to return true for the detailed error")
	}
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("Expected Is to return true for the original error")
	}

	// Test As
	var target *AppError
	if !errors.As(wrappedErr, &target) {
		t.Errorf("Expected As to return true and set target")
	}

	if target != detailedErr {
		t.Errorf("Expected target to be the same as detailed error")
	}
}

func TestAppError_MarshalJSON(t *testing.T) {
	errWithoutDetails := New(Validation, "Validation error", DefaultValidationCode)
	errWithDetails := New(Validation, "Validation error", DefaultValidationCode).WithDetail("field", "username")

	expectedJSONWithoutDetails := `{"message":"Validation error","code":3000,"type":"VALIDATION"}`
	expectedJSONWithDetails := `{"message":"Validation error","code":3000,"type":"VALIDATION","details":{"field":"username"}}`

	actualJSONWithoutDetails, err := json.Marshal(errWithoutDetails)
	if err != nil {
		t.Errorf("Error marshaling AppError to JSON: %v", err)
	}
	actualJSONWithDetails, err := json.Marshal(errWithDetails)
	if err != nil {
		t.Errorf("Error marshaling AppError to JSON: %v", err)
	}

	if string(actualJSONWithoutDetails) != expectedJSONWithoutDetails {
		t.Errorf("Expected marshaled JSON to be %s, got %s", expectedJSONWithoutDetails, string(actualJSONWithoutDetails))
	}
	if string(actualJSONWithDetails) != expectedJSONWithDetails {
		t.Errorf("Expected marshaled JSON to be %s, got %s", expectedJSONWithDetails, string(actualJSONWithDetails))
	}
}
