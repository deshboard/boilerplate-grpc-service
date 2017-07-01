// +build acceptance

package app

import (
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/goph/stdlib/net"
	"google.golang.org/grpc"
)

func FeatureContext(s *godog.Suite) {
	addr := net.ResolveVirtualAddr("pipe", "pipe")
	listener, dialer := net.PipeListen(addr)

	server := grpc.NewServer()
	client, _ := grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(func(s string, t time.Duration) (stdnet.Conn, error) { return dialer.Dial() }))

	// Add steps here

	go server.Serve(listener)
}
