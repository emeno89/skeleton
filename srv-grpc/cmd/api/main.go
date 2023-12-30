package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"skeleton/shared/bootstrap"
	"skeleton/shared/text"
	"skeleton/srv-grpc/internal/api"
	"skeleton/srv-grpc/internal/pkg/example"
)

func main() {
	fx.New(
		fx.Provide(bootstrap.MakeDebugConfigFromEnv),
		fx.Provide(bootstrap.ZapLogger),
		fx.WithLogger(bootstrap.FxLogger),
		fx.Provide(text.NewBundle),
		example.Module,
		api.Module,
	).Run()
}
