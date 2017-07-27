// +build acceptance

package app

import (
	stdnet "net"
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
	client.Close() // Remove this line

	go server.Serve(listener)
}
