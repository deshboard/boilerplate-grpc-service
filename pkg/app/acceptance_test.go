// +build acceptance

package app_test

import (
	stdnet "net"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/deshboard/boilerplate-grpc-service/test"
	"github.com/goph/stdlib/net"
	"google.golang.org/grpc"
)

func init() {
	test.RegisterFeaturePath("../features")
	test.RegisterFeatureContext(FeatureContext)
}

func FeatureContext(s *godog.Suite) {
	addr := net.ResolveVirtualAddr("pipe", "pipe")
	listener, dialer := net.PipeListen(addr)

	server := grpc.NewServer()
	client, _ := grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(func(s string, t time.Duration) (stdnet.Conn, error) { return dialer.Dial() }))

	// Add steps here
	func(s *godog.Suite, server *grpc.Server, client *grpc.ClientConn) {}(s, server, client)

	go server.Serve(listener)
}
