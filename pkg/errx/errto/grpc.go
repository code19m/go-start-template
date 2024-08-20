package errto

import (
	"fmt"
	"go-start-template/pkg/errx"
	"go-start-template/pkg/errx/internal/errpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errto.GRPC converts an error to a gRPC status error.
// This function is intended for use in gRPC server handlers to convert ErrorX instances to gRPC status errors.
func GRPC(err error) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*errx.ErrorX)
	if !ok {
		err = errx.Wrap(err)
	}

	return toStatus(err).Err()
}

func toStatus(err error) *status.Status {
	if e, ok := err.(*errx.ErrorX); ok {
		st, dtErr := status.New(gRPCStatusCode(e), e.Message).WithDetails(
			&errpb.ErrorX{
				Message: e.Message,
				Code:    e.Code,
				Type:    int32(e.Type),
				Details: e.Details,
			},
		)
		if dtErr == nil {
			return st
		}
		err = fmt.Errorf("st.WithDetails error: %w original error: %w", dtErr, err)
	}
	return status.New(codes.Internal, err.Error())
}

func gRPCStatusCode(err error) codes.Code {
	if e, ok := err.(*errx.ErrorX); ok {
		switch e.Type {
		case errx.Authentication:
			return codes.Unauthenticated
		case errx.Forbidden:
			return codes.PermissionDenied
		case errx.Validation:
			return codes.InvalidArgument
		case errx.NotFound:
			return codes.NotFound
		case errx.Conflict:
			return codes.AlreadyExists
		case errx.Internal:
			return codes.Internal
		}
	}
	return codes.Internal
}
