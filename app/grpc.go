package app

import (
	"github.com/goph/fxt/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// NewGrpcConfig creates a grpc config.
func NewGrpcConfig(config Config, tracer opentracing.Tracer) *fxgrpc.Config {
	addr := config.GrpcAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := fxgrpc.NewConfig(addr)
	c.ReflectionEnabled = config.GrpcEnableReflection
	c.Options = []grpc.ServerOption{
		grpc_middleware.WithStreamServerChain(
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_recovery.StreamServerInterceptor(),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	return c
}
