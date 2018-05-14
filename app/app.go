package app

import (
	grpcfx "github.com/deshboard/boilerplate-grpc-service/pkg/driver/fx"
	"github.com/goph/fxt/app/grpc"
	"github.com/goph/fxt/tracing/opentracing"
	"go.uber.org/fx"
)

// Module is the collection of all modules of the application.
var Module = fx.Options(
	fxgrpcapp.Module,

	// Configuration
	fx.Provide(
		NewLoggerConfig,
		NewDebugConfig,
	),

	// gRPC server
	fx.Provide(NewGrpcConfig),

	fx.Provide(fxopentracing.NewTracer),

	grpcfx.Module,
)

type Runner = fxgrpcapp.Runner
