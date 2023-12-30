package endpoints

import (
	"errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNoId = errors.New("no id")

func wrapResponseErr(err error, localizer *i18n.Localizer) error {
	switch err {
	case ErrNoId:
		return status.New(codes.InvalidArgument, ErrNoIdText(localizer)).Err()
	default:
		return status.New(codes.Unknown, err.Error()).Err()
	}
}
