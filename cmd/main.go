package main

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/errors"
	"github.com/goph/fxt/grpc"
	fxlog "github.com/goph/fxt/log"
	"github.com/goph/fxt/tracing"
	"github.com/goph/healthz"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
)

func main() {
	status := healthz.NewStatusChecker(healthz.Healthy)
	var ext struct {
		Config       *Config
		Closer       fxt.Closer
		Logger       log.Logger
		ErrorHandler emperror.Handler

		DebugErr debug.Err
		GrpcErr  grpc.Err
	}

	app := fx.New(
		fx.NopLogger,
		fxt.Bootstrap,
		fx.Provide(
			NewConfig,

			// Log and error handling
			NewLoggerConfig,
			fxlog.NewLogger,
			errors.NewHandler,

			// Debug server
			NewDebugConfig,
			debug.NewServer,
			debug.NewHealthCollector,
		),
		fx.Invoke(grpc_prometheus.Register, RegisterPrometheusHandler),
		fx.Invoke(func(collector healthz.Collector) {
			collector.RegisterChecker(healthz.ReadinessCheck, status)
		}),
		fx.Extract(&ext),

		fx.Provide(
			// gRPC server
			NewService,
			NewGrpcConfig,
			NewStreamInterceptor,
			NewUnaryInterceptor,
			grpc.NewServer,

			tracing.NewTracer,
		),
	)

	err := app.Err()
	if err != nil {
		panic(err)
	}

	// Close resources even when there is an error
	defer ext.Closer.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(ext.ErrorHandler)

	level.Info(ext.Logger).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
	)

	err = app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	select {
	case sig := <-app.Done():
		level.Info(ext.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-ext.DebugErr:
		if err != nil {
			err = emperror.WithStack(emperror.WithMessage(err, "debug server crashed"))
			ext.ErrorHandler.Handle(err)
		}

	case err := <-ext.GrpcErr:
		if err != nil {
			err = emperror.WithStack(emperror.WithMessage(err, "grpc server crashed"))
			ext.ErrorHandler.Handle(err)
		}
	}

	status.SetStatus(healthz.Unhealthy)

	ctx, cancel := context.WithTimeout(context.Background(), ext.Config.ShutdownTimeout)
	defer cancel()

	err = app.Stop(ctx)
	emperror.HandleIfErr(ext.ErrorHandler, err)
}
