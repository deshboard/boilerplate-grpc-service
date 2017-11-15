package acceptance

import (
	"net"

	"github.com/DATA-DOG/godog"
	"google.golang.org/grpc"
)

type GrpcFeatureContext struct {
	server *grpc.Server
	clientConn *grpc.ClientConn

	frozen bool
}

func (c *GrpcFeatureContext) FeatureContext(s *godog.Suite) {
	if c.frozen {
		panic("trying to use a frozen feature context")
	}
	c.frozen = true

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	clientConn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	s.BeforeSuite(func() { go server.Serve(lis) })
	s.AfterSuite(server.Stop)

	c.server = server
	c.clientConn = clientConn
}

func (c *GrpcFeatureContext) Server() *grpc.Server {
	return c.server
}

func (c *GrpcFeatureContext) ClientConn() *grpc.ClientConn {
	return c.clientConn
}
