package main // import "github.com/deshboard/boilerplate-grpc-service"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-grpc-service/app"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
	"github.com/sagikazarmark/utilz/util"
	"google.golang.org/grpc"
)

func main() {
	defer logger.Info("Shutting down")
	defer shutdownManager.Shutdown()

	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
	}).Infof("Starting %s", app.FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdownManager.Register(w.Close)

	serverManager := serverz.NewServerManager(logger)
	errChan := make(chan error, 10)
	signalChan := make(chan os.Signal, 1)

	var debugServer serverz.Server
	if config.Debug {
		debugServer = &serverz.NamedServer{
			Server: &http.Server{
				Handler:  http.DefaultServeMux,
				ErrorLog: log.New(w, "debug: ", 0),
			},
			Name: "debug",
		}

		shutdownManager.RegisterAsFirst(debugServer.Close)
		go serverManager.ListenAndStartServer(debugServer, config.DebugAddr)(errChan)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_logrus.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	grpc_prometheus.Register(grpcServer)

	grpcServerWrapper := &serverz.NamedServer{
		Server: &serverz.GrpcServer{grpcServer},
		Name:   "grpc",
	}

	serviceHealth := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceHealth)

	status := healthz.NewStatusChecker(healthz.Healthy)
	checkerCollector.RegisterChecker(healthz.ReadinessCheck, status)
	healthService := checkerCollector.NewHealthService()
	healthHandler := http.NewServeMux()

	healthHandler.Handle("/healthz", healthService.Handler(healthz.LivenessCheck))
	healthHandler.Handle("/readiness", healthService.Handler(healthz.ReadinessCheck))

	if config.MetricsEnabled {
		logger.Debug("Serving metrics under health endpoint")

		healthHandler.Handle("/metrics", promhttp.Handler())
	}

	healthServer := &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: log.New(w, "health: ", 0),
		},
		Name: "health",
	}

	shutdownManager.RegisterAsFirst(util.ShutdownFunc(grpcServer.Stop), healthServer.Close)
	go serverManager.ListenAndStartServer(grpcServerWrapper, config.ServiceAddr)(errChan)
	go serverManager.ListenAndStartServer(healthServer, config.HealthAddr)(errChan)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			status.SetStatus(healthz.Unhealthy)

			if err != nil {
				logger.Error(err)
			} else {
				logger.Warning("Error channel received non-error value")
			}

			// Break the loop, proceed with regular shutdown
			break MainLoop
		case s := <-signalChan:
			logger.Infof(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			logger.Debugf("Shutting down with timeout %v", config.ShutdownTimeout)

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			wg := &sync.WaitGroup{}

			if config.Debug {
				go serverManager.StopServer(debugServer, wg)(ctx)
			}

			go serverManager.StopServer(grpcServerWrapper, wg)(ctx)
			go serverManager.StopServer(healthServer, wg)(ctx)

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}
}
