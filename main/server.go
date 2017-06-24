package main

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	_grpc "github.com/goph/serverz/grpc"
	"github.com/goph/stdlib/ext"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	opentracing "github.com/opentracing/opentracing-go"
)

// newServer creates the main server instance for the service
func newServer(config *configuration, logger log.Logger, tracer opentracing.Tracer, healthCollector healthz.Collector) (serverz.Server, ext.Closer) {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	server := createGrpcServer(tracer)

	// Register servers here

	grpc_prometheus.Register(server)

	return &serverz.NamedServer{
		Server: &_grpc.Server{server},
		Name:   "grpc",
	}, ext.NoopCloser
}
