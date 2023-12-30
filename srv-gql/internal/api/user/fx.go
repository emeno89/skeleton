package user

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			newJwtJen,
			fx.As(new(jwtGen)),
		),
	),
	fx.Provide(newExampleLister),
	fx.Provide(
		fx.Annotate(
			newUserLister,
			fx.As(new(userLister)),
		),
	),
	fx.Provide(newService),
)
