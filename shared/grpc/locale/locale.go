package locale

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

const acceptLanguageHeader = "Accept-Language"

func AddAcceptLanguage(ctx context.Context, locale string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, strings.ToLower(acceptLanguageHeader), locale)
}

func GetAcceptLanguages(ctx context.Context) []string {
	return metadata.ValueFromIncomingContext(ctx, strings.ToLower(acceptLanguageHeader))
}
