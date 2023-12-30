package example

import (
	"context"
	"skeleton/shared/grpc/dial"
	"skeleton/shared/grpc/pb/grpcpb"
	"time"
)

const (
	baseReqTimeout = 10 * time.Second
)

type Client interface {
	GetManyById(ctx context.Context, ids ...string) ([]*grpcpb.ExampleItem, error)
}

type grpcClient struct {
	conn *dial.Conn
}

func newGrpcClient(cfg clientConfig) *grpcClient {
	return &grpcClient{
		conn: dial.NewConn(cfg.Host),
	}
}

func (s *grpcClient) GetManyById(ctx context.Context, ids ...string) ([]*grpcpb.ExampleItem, error) {
	conn, err := s.conn.Open(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, baseReqTimeout)
	defer cancel()

	req := &grpcpb.GetManyByIdRequest{Ids: ids}

	resp, err := grpcpb.NewExampleServiceClient(conn).GetManyById(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}
