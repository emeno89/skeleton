package api

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"skeleton/shared/bootstrap"
	"skeleton/srv-gql/internal/api/graphql"
	"skeleton/srv-gql/internal/api/user"
)

var Module = fx.Options(
	fx.Provide(makeConfigFromEnv),
	user.Module,
	graphql.Module,
	fx.Provide(newServer),
	fx.Invoke(func(lc fx.Lifecycle, srv *server, logger *zap.Logger) {
		logger = logger.With(zap.String("addr", srv.httpSrv.Addr))

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return fxStartEvent(ctx, srv, logger)
			},
			OnStop: func(ctx context.Context) error {
				return fxStopEvent(ctx, srv, logger)
			},
		})
	}),
	fx.Invoke(func(cfg serverConfig, logger *zap.Logger) {
		bootstrap.HealthCheck(cfg.HealthPort, logger)
	}),
)

func fxStartEvent(_ context.Context, srv *server, logger *zap.Logger) error {
	errCh := make(chan error)

	go func() {
		logger.Info("starting api server")

		errCh <- srv.start()
	}()

	go func() {
		if err := <-errCh; err != http.ErrServerClosed {
			logger.Fatal("api server failed", zap.Error(err))
		}
	}()

	return nil
}

func fxStopEvent(ctx context.Context, srv *server, logger *zap.Logger) error {
	logger.Info("stopping api server")

	return srv.stop(ctx)
}
