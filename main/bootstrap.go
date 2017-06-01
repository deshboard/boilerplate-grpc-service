package main

import (
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func bootstrap() serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	server := createGrpcServer()

	// Register servers here

	grpc_prometheus.Register(server)

	return &serverz.NamedServer{
		Server: &serverz.GrpcServer{server},
		Name:   "grpc",
	}
}
