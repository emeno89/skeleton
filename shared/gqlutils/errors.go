package gqlutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CodeUnauthorized    = "Unauthorized"
	CodeNotFound        = "NotFound"
	CodeForbidden       = "Forbidden"
	CodeValidationError = "InvalidArgument"
)

type ApiError struct {
	gqlErr *gqlerror.Error
}

func NewApiError(ctx context.Context, text string) *ApiError {
	return &ApiError{
		gqlErr: gqlerror.ErrorPathf(graphql.GetFieldContext(ctx).Path(), text),
	}
}

func (e *ApiError) Error() string {
	return e.gqlErr.Error()
}

func (e *ApiError) GetGQL() *gqlerror.Error {
	return e.gqlErr
}

func (e *ApiError) init() {
	if e.gqlErr.Extensions == nil {
		e.gqlErr.Extensions = make(map[string]interface{})
	}
}

func (e *ApiError) SetCode(code string) *ApiError {
	e.init()
	e.gqlErr.Extensions["code"] = code

	return e
}

func (e *ApiError) SetField(field string) *ApiError {
	e.init()
	e.gqlErr.Extensions["field"] = field

	return e
}

func ErrNotAuthorized(ctx context.Context) error {
	err := NewApiError(ctx, "Not Authorized")

	_ = err.SetCode(CodeUnauthorized)

	return err
}

func ErrNotFound(ctx context.Context) error {
	err := NewApiError(ctx, "Not Found")

	_ = err.SetCode(CodeNotFound)

	return err
}

func ErrUnhandled(ctx context.Context) *gqlerror.Error {
	fieldCtx := graphql.GetFieldContext(ctx)

	return gqlerror.ErrorPathf(fieldCtx.Path(), "Unhandled error")
}

func ErrForbidden(ctx context.Context) error {
	err := NewApiError(ctx, "Forbidden")

	_ = err.SetCode(CodeForbidden)

	return err
}

func ConvertGPRCError(ctx context.Context, err error) error {
	gStatus, _ := status.FromError(err)

	if isUnexpectedCode(gStatus) {
		return err
	}

	if gStatus.Code() == codes.InvalidArgument {
		return convertInvalidArgumentErr(ctx, gStatus)
	}

	return NewApiError(ctx, gStatus.Message())
}

func convertInvalidArgumentErr(ctx context.Context, gStatus *status.Status) error {
	details := gStatus.Details()
	if len(details) == 0 {
		return NewApiError(ctx, gStatus.Message())
	}

	converted, ok := details[0].(*errdetails.BadRequest)
	if !ok {
		return NewApiError(ctx, gStatus.Message())
	}

	violations := converted.GetFieldViolations()
	if len(violations) == 0 {
		return NewApiError(ctx, gStatus.Message())
	}

	var errs []error

	for _, v := range violations {
		errs = append(errs, NewApiError(ctx, v.Description).SetField(v.Field).SetCode(CodeValidationError))
	}

	for _, err := range errs[1:] {
		graphql.AddError(ctx, err)
	}

	return errs[0]
}

func isUnexpectedCode(gStatus *status.Status) bool {
	unexpectedCodes := []codes.Code{
		codes.Canceled,
		codes.Unknown,
		codes.DeadlineExceeded,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss,
		codes.Unauthenticated,
	}

	for _, code := range unexpectedCodes {
		if code == gStatus.Code() {
			return true
		}
	}

	return false
}

func PresentErr(ctx context.Context, err error, logger *zap.Logger) *gqlerror.Error {
	var apiErrPtr *ApiError

	if errors.As(err, &apiErrPtr) {
		return apiErrPtr.GetGQL()
	}

	oCtx := graphql.GetOperationContext(ctx)
	fCtx := graphql.GetFieldContext(ctx)

	logger.Error(
		"Unhandled error",
		zap.Error(err),
		zap.Any("path", fCtx.Path()),
		zap.String("raw", oCtx.RawQuery),
		zap.Any("vars", oCtx.Variables),
	)

	return ErrUnhandled(ctx)
}

func RecoverErr(ctx context.Context, err interface{}, logger *zap.Logger) error {
	op := graphql.GetOperationContext(ctx)

	logger.Error("Internal server error", zap.Any("err", err), zap.String("query", op.RawQuery))

	return fmt.Errorf("internal server error")
}
