package endpoints

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"skeleton/shared/grpc/locale"
	"skeleton/shared/grpc/pb/grpcpb"
	"skeleton/srv-grpc/internal/pkg/example"
)

type Endpoints struct {
	grpcpb.UnimplementedExampleServiceServer
	bundle         *i18n.Bundle
	exampleService *example.Service
}

func newEndpoints(bundle *i18n.Bundle, exampleService *example.Service) *Endpoints {
	return &Endpoints{
		bundle:         bundle,
		exampleService: exampleService,
	}
}

func (e *Endpoints) GetManyById(ctx context.Context, req *grpcpb.GetManyByIdRequest) (*grpcpb.ManyByIdResponse, error) {
	localizer := i18n.NewLocalizer(e.bundle, locale.GetAcceptLanguages(ctx)...)

	if len(req.Ids) == 0 {
		return nil, wrapResponseErr(ErrNoId, localizer)
	}

	result, err := e.exampleService.GetManyId(ctx, req.Ids...)
	if err != nil {
		return nil, wrapResponseErr(err, localizer)
	}

	resp := &grpcpb.ManyByIdResponse{
		Items: make([]*grpcpb.ExampleItem, len(result)),
	}

	for i, val := range result {
		resp.Items[i] = newExampleItem(val)
	}

	return resp, nil
}
