package graph

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"net/http"
	"skeleton/srv-gql/internal/api/graphql/graph/loader"
	"skeleton/srv-gql/internal/api/user"
	"skeleton/srv-gql/internal/pkg/client/example"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	bundle        *i18n.Bundle
	userService   *user.Service
	exampleClient example.Client
	logger        *zap.Logger
}

func newResolver(
	bundle *i18n.Bundle,
	userService *user.Service,
	exampleClient example.Client,
	logger *zap.Logger,
) *Resolver {
	return &Resolver{
		bundle:        bundle,
		userService:   userService,
		exampleClient: exampleClient,
		logger:        logger,
	}
}

func (r *Resolver) LoaderHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			ctx := loader.ToOutgoingContext(req.Context(), r.exampleClient)

			next.ServeHTTP(w, req.WithContext(ctx))
		},
	)
}
