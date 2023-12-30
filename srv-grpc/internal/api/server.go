package api

import (
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"net"
	"skeleton/shared/grpc/pb/grpcpb"
	"skeleton/srv-grpc/internal/api/endpoints"
)

type server struct {
	grpcSrv   *grpc.Server
	endpoints *endpoints.Endpoints
}

func newServer(endpoints *endpoints.Endpoints) *server {
	return &server{
		grpcSrv: grpc.NewServer(
			grpc.UnaryInterceptor(
				apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			),
		),
		endpoints: endpoints,
	}
}

func (s *server) start(cfg serverConfig) error {
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return err
	}

	grpcpb.RegisterExampleServiceServer(s.grpcSrv, s.endpoints)

	return s.grpcSrv.Serve(listener)
}

func (s *server) stop() {
	s.grpcSrv.GracefulStop()
}
