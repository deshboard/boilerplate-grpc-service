// +build acceptance

package app_test

import (
	"net"

	"github.com/DATA-DOG/godog"
	"github.com/deshboard/boilerplate-grpc-service/test"
	"google.golang.org/grpc"
)

func init() {
	test.RegisterFeaturePath("../features")
	test.RegisterFeatureContext(FeatureContext)
}

func FeatureContext(s *godog.Suite) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	client, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	// Add steps here
	func(s *godog.Suite, server *grpc.Server, client *grpc.ClientConn) {}(s, server, client)

	go server.Serve(lis)
}
