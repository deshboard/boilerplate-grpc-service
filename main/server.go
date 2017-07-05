package main

import (
	"time"

	"github.com/goph/healthz"
	"github.com/goph/serverz"
	_grpc "github.com/goph/serverz/grpc"
	"github.com/goph/serverz/named"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// newServer creates the main server instance for the service.
func newServer(appCtx *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	server := createGrpcServer(appCtx)

	// Register servers here

	grpc_prometheus.Register(server)

	return &named.Server{
		Server:     &_grpc.Server{server},
		ServerName: "grpc",
	}
}
