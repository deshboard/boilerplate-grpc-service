// +build acceptance

package app

import (
	"os"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/goph/stdlib/net"
	"google.golang.org/grpc"
)

func init() {
	runs = append(runs, func() int {
		format := "progress"
		for _, arg := range os.Args[1:] {
			// go test transforms -v option
			if arg == "-test.v=true" {
				format = "pretty"
				break
			}
		}

		return godog.RunWithOptions(
			"godog",
			FeatureContext,
			godog.Options{
				Format:    format,
				Paths:     []string{"features"},
				Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
			},
		)
	})
}

func FeatureContext(s *godog.Suite) {
	addr := net.ResolveVirtualAddr("pipe", "pipe")
	listener, dialer := net.PipeListen(addr)

	server := grpc.NewServer()
	client, _ := grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(func(s string, t time.Duration) (stdnet.Conn, error) { return dialer.Dial() }))

	// Add steps here

	go server.Serve(listener)
}
