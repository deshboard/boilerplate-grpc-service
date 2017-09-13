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
	"google.golang.org/grpc/reflection"
)

func createGrpcServer(app *application) *grpc.Server {
	logger := log.With(
		app.logger,
		"server", "grpc",
	)

	// TODO: separate log levels
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		log.NewStdlibAdapter(level.Info(logger)),
		log.NewStdlibAdapter(level.Warn(logger)),
		log.NewStdlibAdapter(level.Error(logger)),
	))

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(app.tracer)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(app.tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	if app.config.GrpcEnableReflection {
		level.Debug(app.logger).Log("msg", "grpc reflection enabled")

		reflection.Register(server)
	}

	return server
}
