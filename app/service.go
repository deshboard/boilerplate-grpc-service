package app

import (
	"github.com/deshboard/boilerplate-grpc-service/pkg/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	fxgrpc "github.com/goph/fxt/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *app.Service {
	return app.NewService(
		app.Logger(params.Logger),
		app.ErrorHandler(params.ErrorHandler),
	)
}

// NewGrpcConfig creates a grpc config.
func NewGrpcConfig(config *Config, tracer opentracing.Tracer) *fxgrpc.Config {
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
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		),
		grpc_middleware.WithUnaryServerChain(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	return c
}
