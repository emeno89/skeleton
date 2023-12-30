package endpoints

import (
	"skeleton/shared/grpc/pb/grpcpb"
	"skeleton/srv-grpc/internal/pkg/example"
)

func newExampleItem(data example.Data) *grpcpb.ExampleItem {
	return &grpcpb.ExampleItem{
		Id: data.Id,
	}
}
