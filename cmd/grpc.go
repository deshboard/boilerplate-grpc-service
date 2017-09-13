package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func createGrpcServer(appCtx *application) *grpc.Server {
	// TODO: separate log levels
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		log.NewStdlibAdapter(level.Info(appCtx.logger)),
		log.NewStdlibAdapter(level.Warn(appCtx.logger)),
		log.NewStdlibAdapter(level.Error(appCtx.logger)),
	))

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
