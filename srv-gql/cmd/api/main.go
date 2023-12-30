package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"skeleton/shared/bootstrap"
	"skeleton/shared/text"
	"skeleton/srv-gql/internal/api"
	"skeleton/srv-gql/internal/pkg/client/example"
	"skeleton/srv-gql/internal/pkg/security"
)

func main() {
	fx.New(
		fx.Provide(bootstrap.MakeDebugConfigFromEnv),
		fx.Provide(bootstrap.ZapLogger),
		fx.WithLogger(bootstrap.FxLogger),
		fx.Provide(text.NewBundle),
		security.Module,
		example.Module,
		api.Module,
	).Run()
}
