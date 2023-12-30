package dial

import (
	"context"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

const dialTimeout = 10 * time.Second

type Conn struct {
	grpcConn *grpc.ClientConn
	host     string
	mux      *sync.Mutex
}

func NewConn(host string) *Conn {
	return &Conn{
		host: host,
		mux:  new(sync.Mutex),
	}
}

func (s *Conn) Open(ctx context.Context) (*grpc.ClientConn, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.grpcConn != nil {
		return s.grpcConn, nil
	}

	ctx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()

	cc, err := grpc.DialContext(
		ctx,
		s.host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
	)

	if err != nil {
		return nil, err
	}

	s.grpcConn = cc

	return s.grpcConn, nil
}
