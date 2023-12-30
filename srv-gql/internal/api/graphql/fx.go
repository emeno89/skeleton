package graphql

import (
	"go.uber.org/fx"
	"skeleton/srv-gql/internal/api/graphql/graph"
)

var Module = fx.Options(
	graph.Module,
	fx.Provide(newRouter),
)
