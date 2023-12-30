package locale

import (
	"context"
	"net/http"
)

type dataKey struct{}

func FromIncomingContext(ctx context.Context) string {
	locale, _ := ctx.Value(dataKey{}).(string)

	return locale
}

func Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dataKey{}, r.Header.Get("Accept-Language"))

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
