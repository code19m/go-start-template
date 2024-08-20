package errx

import (
	"errors"
	"go-start-template/pkg/errx/internal/errpb"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fromPG(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505", "21000":
			return ErrConflict.
				WithDetail("constraint", pgErr.ConstraintName)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func fromGRPC(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, detail := range st.Details() {
		if pb, ok := detail.(*errpb.ErrorX); ok {
			return fromProto(pb)
		}
	}

	return fromGRPCCode(st)
}

func fromGRPCCode(st *status.Status) *ErrorX {
	msg := st.Message()

	grpcToAppErr := map[codes.Code]*ErrorX{
		codes.AlreadyExists:    ErrConflict,
		codes.NotFound:         ErrNotFound,
		codes.PermissionDenied: ErrForbidden,
		codes.Unauthenticated:  ErrAuthentication,
	}

	if st.Code() == codes.InvalidArgument {
		return handleInvalidArgument(st, msg)
	}

	if e, found := grpcToAppErr[st.Code()]; found {
		return e.WithDetail("error", msg)
	}

	return ErrInternal.WithDetail("error", msg)
}

func handleInvalidArgument(st *status.Status, msg string) *ErrorX {
	err := ErrValidation.WithDetail("error", msg)
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok {
			for _, violation := range badRequest.GetFieldViolations() {
				err = err.WithDetail(violation.GetField(), violation.GetDescription())
			}
		}
	}
	return err
}

func fromProto(pberr *errpb.ErrorX) *ErrorX {
	err := New(
		Type(pberr.GetType()),
		pberr.Message,
		pberr.Code,
	)
	for k, v := range pberr.Details {
		err = err.WithDetail(k, v)
	}
	return err
}
