package main

import (
	"time"

	"github.com/goph/healthz"
	"github.com/goph/serverz/aio"
	"github.com/goph/serverz/grpc"
	"github.com/goph/stdlib/net"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// newGrpcServer creates the main server instance for the service.
func newGrpcServer(appCtx *application) *aio.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.GrpcAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	server := createGrpcServer(appCtx)

	// Register servers here

	grpc_prometheus.Register(server)

	return &aio.Server{
		Server: &grpc.Server{Server: server},
		Name:   "grpc",
		Addr:   net.ResolveVirtualAddr("tcp", appCtx.config.GrpcAddr),
	}
}
