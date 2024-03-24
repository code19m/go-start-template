package apperr_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-start-template/pkg/apperr"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	t.Parallel()

	var err error = apperr.New(apperr.Authorization, "Unauthorized", apperr.DefaultAuthorizationCode)

	appErr, ok := err.(*apperr.AppError)
	if !ok {
		t.Errorf("Expected AppError, got %T", err)
	}

	if appErr.Message != "Unauthorized" {
		t.Errorf("Expected message 'Unauthorized', got %s", appErr.Message)
	}

	if appErr.Code != apperr.DefaultAuthorizationCode {
		t.Errorf("Expected code %d, got %d", apperr.DefaultAuthorizationCode, appErr.Code)
	}

	if appErr.Type != apperr.Authorization {
		t.Errorf("Expected type Authorization, got %s", appErr.Type)
	}

	if appErr.Error() != "Unauthorized" {
		t.Errorf("Expected error message 'Unauthorized', got %s", appErr.Error())
	}
}

func TestWithDetail(t *testing.T) {
	t.Parallel()

	err := apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode)
	err = err.WithDetail("field", "username")

	if len(err.Details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(err.Details))
	}

	if err.Details["field"] != "username" {
		t.Errorf("Expected detail 'username' for 'field', got %s", err.Details["field"])
	}
}

func TestHTTPCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err    error
		status int
	}{
		{apperr.New(apperr.Authorization, "Unauthorized", apperr.DefaultAuthorizationCode), http.StatusUnauthorized},
		{apperr.New(apperr.Forbidden, "Forbidden", apperr.DefaultForbiddenCode), http.StatusForbidden},
		{apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode), http.StatusBadRequest},
		{apperr.New(apperr.NotFound, "Not found", apperr.DefaultNotFoundCode), http.StatusNotFound},
		{apperr.New(apperr.Conflict, "Conflict", apperr.DefaultConflictCode), http.StatusConflict},
		{apperr.New(apperr.Internal, "Internal error", apperr.DefaultInternalCode), http.StatusInternalServerError},
		{errors.New("unknown error"), http.StatusInternalServerError},
	}

	for _, test := range tests {
		if apperr.HTTPCode(test.err) != test.status {
			t.Errorf("Expected HTTP status %d for error, got %d", test.status, apperr.HTTPCode(test.err))
		}
	}
}

func TestGRPCCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		code int
	}{
		{apperr.New(apperr.Authorization, "Unauthorized", apperr.DefaultAuthorizationCode), int(codes.Unauthenticated)},
		{apperr.New(apperr.Forbidden, "Forbidden", apperr.DefaultForbiddenCode), int(codes.PermissionDenied)},
		{apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode), int(codes.InvalidArgument)},
		{apperr.New(apperr.NotFound, "Not found", apperr.DefaultNotFoundCode), int(codes.NotFound)},
		{apperr.New(apperr.Conflict, "Conflict", apperr.DefaultConflictCode), int(codes.AlreadyExists)},
		{apperr.New(apperr.Internal, "Internal error", apperr.DefaultInternalCode), int(codes.Internal)},
		{errors.New("unknown error"), int(codes.Internal)},
	}

	for _, test := range tests {
		if apperr.GRPCCode(test.err) != test.code {
			t.Errorf("Expected gRPC code %d for error, got %d", test.code, apperr.GRPCCode(test.err))
		}
	}
}

func TestUnwrapIsAs(t *testing.T) {
	t.Parallel()

	originalErr := apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode)
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
	var target *apperr.AppError
	if !errors.As(wrappedErr, &target) {
		t.Errorf("Expected As to return true and set target")
	}

	if target != originalErr {
		t.Errorf("Expected target to be the same as appErr")
	}
}

func TestUnwrapIsAsWithDetails(t *testing.T) {
	t.Parallel()

	originalErr := apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode)
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
	var target *apperr.AppError
	if !errors.As(wrappedErr, &target) {
		t.Errorf("Expected As to return true and set target")
	}

	if target != detailedErr {
		t.Errorf("Expected target to be the same as detailed error")
	}
}

func TestAppError_MarshalJSON(t *testing.T) {
	t.Parallel()

	errWithoutDetails := apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode)
	errWithDetails := apperr.New(apperr.Validation, "Validation error", apperr.DefaultValidationCode).WithDetail("field", "username")

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
