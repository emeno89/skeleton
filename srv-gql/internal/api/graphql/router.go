package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
	"skeleton/shared/gqlutils"
	"skeleton/srv-gql/internal/api/graphql/graph"
	"skeleton/srv-gql/internal/api/graphql/graph/generated"
	"skeleton/srv-gql/internal/api/user"
)

type Router struct {
	resolver    *graph.Resolver
	userService *user.Service
	logger      *zap.Logger
}

func newRouter(
	resolver *graph.Resolver,
	userService *user.Service,
	logger *zap.Logger,
) *Router {
	return &Router{
		userService: userService,
		resolver:    resolver,
		logger:      logger,
	}
}

func (s *Router) Mount(router chi.Router, mountPath string) {
	cfg := generated.Config{Resolvers: s.resolver}

	s.setDirectives(&cfg)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	srv.SetRecoverFunc(s.recoverFunc)
	srv.SetErrorPresenter(s.errorPresenterFunc)
	srv.AroundFields(s.aroundFieldsFunc)

	h := s.userService.Handler(
		s.resolver.LoaderHandler(srv),
	)

	router.Handle("/graphiql", playground.AltairHandler("GraphQL playground", mountPath+"/graphql"))
	router.Handle("/graphql", h)
}

func (s *Router) recoverFunc(ctx context.Context, err interface{}) error {
	return gqlutils.RecoverErr(ctx, err, s.logger)
}

func (s *Router) errorPresenterFunc(ctx context.Context, err error) *gqlerror.Error {
	return gqlutils.PresentErr(ctx, err, s.logger)
}

func (s *Router) setDirectives(cfg *generated.Config) {
	cfg.Directives.IsLoggedIn = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		if _, ok := s.userService.FromIncomingContext(ctx); !ok {
			return nil, gqlutils.ErrNotAuthorized(ctx)
		}

		return next(ctx)
	}
}

func (s *Router) aroundFieldsFunc(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	return gqlutils.ApmTraceField(ctx, next)
}
