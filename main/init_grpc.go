package main

import "google.golang.org/grpc/grpclog"

func init() {
	// Set global gRPC logger
	grpclog.SetLogger(logger.WithField("server", "grpc"))
}
