package api

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"skeleton/shared/bootstrap"
	"skeleton/srv-grpc/internal/api/endpoints"
)

var Module = fx.Options(
	endpoints.Module,
	fx.Provide(makeConfigFromEnv),
	fx.Provide(newServer),
	fx.Invoke(func(lc fx.Lifecycle, cfg serverConfig, srv *server, logger *zap.Logger) {
		logger = logger.With(zap.String("addr", ":"+cfg.Port))

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return fxStartEvent(cfg, srv, logger)
			},
			OnStop: func(ctx context.Context) error {
				return fxStopEvent(srv, logger)
			},
		})
	}),
	fx.Invoke(func(cfg serverConfig, logger *zap.Logger) {
		bootstrap.HealthCheck(cfg.HealthPort, logger)
	}),
)

func fxStartEvent(cfg serverConfig, srv *server, logger *zap.Logger) error {
	errCh := make(chan error)

	go func() {
		logger.Info("starting api server")

		errCh <- srv.start(cfg)
	}()

	go func() {
		if err := <-errCh; err != nil {
			logger.Fatal("api server failed", zap.Error(err))
		}
	}()

	return nil
}

func fxStopEvent(srv *server, logger *zap.Logger) error {
	logger.Info("stopping api server")

	srv.stop()

	return nil
}
