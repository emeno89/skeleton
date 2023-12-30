package example

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(makeConfigFromEnv),
	fx.Provide(
		fx.Annotate(
			newGrpcClient,
			fx.As(new(Client)),
		),
	),
)
