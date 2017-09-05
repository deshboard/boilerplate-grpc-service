package main

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func init() {
	// TODO: Set global gRPC logger
	// grpclog.SetLogger()
}

func createGrpcServer(appCtx *application) *grpc.Server {
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(appCtx.tracer)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(appCtx.tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	return server
}
