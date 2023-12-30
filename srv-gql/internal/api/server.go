package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.elastic.co/apm/module/apmchiv5"
	"net/http"
	"skeleton/srv-gql/internal/api/graphql"
	"skeleton/srv-gql/internal/api/locale"
)

const mountPath = "/api/v1"

type routerMount interface {
	Mount(router chi.Router, mountPath string)
}

// @title						Gql server API
// @version					1.0
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
// @BasePath					/api/v1/graphql
type server struct {
	cfg     serverConfig
	httpSrv *http.Server
	mounts  []routerMount
}

func newServer(
	cfg serverConfig,
	graphRouter *graphql.Router,
) *server {
	return &server{
		cfg: cfg,
		httpSrv: &http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Port),
		},
		mounts: []routerMount{graphRouter},
	}
}

func (s *server) start() error {
	router := chi.NewRouter()

	router.Use(locale.Handler)
	router.Use(cors.AllowAll().Handler)
	router.Use(apmchiv5.Middleware())

	if s.cfg.SwaggerEnabled == 1 {
		router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(mountPath+"/swagger/doc.json")))
	}

	for _, mnt := range s.mounts {
		mnt.Mount(router, mountPath)
	}

	router.Mount(mountPath, router)

	s.httpSrv.Handler = router

	return s.httpSrv.ListenAndServe()
}

func (s *server) stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
