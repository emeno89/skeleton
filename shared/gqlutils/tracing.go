package gqlutils

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"go.elastic.co/apm"
	"strings"
)

const (
	traceType = "graphql"
)

func ApmTraceField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fieldCtx := graphql.GetFieldContext(ctx)

	if !fieldCtx.IsResolver {
		return next(ctx)
	}

	objName := strings.ToLower(fieldCtx.Object)
	pathName := fieldCtx.Path().String()

	span, ctx := apm.StartSpan(ctx, fmt.Sprintf("%s %s", objName, pathName), fmt.Sprintf("%s.%s", traceType, objName))

	res, err = next(ctx)

	span.End()

	return res, err
}
