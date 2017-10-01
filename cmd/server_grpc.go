package main

import (
	"time"

	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/serverz/grpc"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
)

// newGrpcServer creates the main server instance for the service.
func newGrpcServer(app *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(app.config.GrpcAddr, healthz.WithTCPTimeout(2*time.Second))
	app.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	server := createGrpcServer(app)

	// Register servers here

	grpc_prometheus.Register(server)

	return &serverz.AppServer{
		Server: &grpc.Server{Server: server},
		Name:   "grpc",
		Addr:   serverz.NewAddr("tcp", app.config.GrpcAddr),
		Logger: app.Logger(),
	}
}
