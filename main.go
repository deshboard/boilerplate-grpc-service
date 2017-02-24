package main // import "github.com/deshboard/boilerplate-grpc-service"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-grpc-service/app"
	"github.com/deshboard/boilerplate-grpc-service/model/boilerplate"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sagikazarmark/healthz"
)

// Global context variables
var (
	config   = &app.Configuration{}
	logger   = logrus.New().WithField("service", app.ServiceName) // Use logrus.FieldLogger type
	tracer   = opentracing.GlobalTracer()
	shutdown = []shutdownHandler{}
)

func main() {
	defer handleShutdown()

	var (
		serviceAddr = flag.String("service", "0.0.0.0:80", "gRPC service address.")
		healthAddr  = flag.String("health", "0.0.0.0:90", "Health service address.")
	)
	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
	}).Printf("Starting %s service", app.FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdown = append(shutdown, w.Close)

	grpcServer := grpc.NewServer()
	boilerplate.RegisterBoilerplateServer(grpcServer, app.NewService())

	healthHandler, status := newHealthServiceHandler()
	healthServer := &http.Server{
		Addr:     *healthAddr,
		Handler:  healthHandler,
		ErrorLog: log.New(w, fmt.Sprintf("%s Health service: ", app.FriendlyServiceName), 0),
	}

	// Force closing server connections (if graceful closing fails)
	shutdown = append([]shutdownHandler{shutdownFunc(grpcServer.Stop), healthServer.Close}, shutdown...)

	errChan := make(chan error, 10)

	go func() {
		logger.WithField("addr", healthServer.Addr).Infof("%s Health service started", app.FriendlyServiceName)
		errChan <- healthServer.ListenAndServe()
	}()

	go func() {
		lis, err := net.Listen("tcp", *serviceAddr)
		if err != nil {
			errChan <- err
			return
		}

		logger.WithField("addr", lis.Addr()).Infof("%s service started", app.FriendlyServiceName)
		errChan <- grpcServer.Serve(lis)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

MainLoop:
	for {
		select {
		case err := <-errChan:
			// In theory this can only be non-nil
			if err != nil {
				// This will be handled (logged) by shutdown
				panic(err)
			} else {
				logger.Info("Error channel received non-error value")

				// Break the loop, proceed with regular shutdown
				break MainLoop
			}
		case s := <-signalChan:
			logger.Println(fmt.Sprintf("Captured %v", s))
			status.SetStatus(healthz.Unhealthy)

			shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				// TODO: implement timeout
				grpcServer.GracefulStop()

				wg.Done()
			}()

			go func() {
				err := healthServer.Shutdown(shutdownContext)
				if err != nil {
					logger.Error(err)
				}

				wg.Done()
			}()

			wg.Wait()

			// Cancel context if shutdown completed earlier
			shutdownCancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}

	close(errChan)
	close(signalChan)
}

type shutdownHandler func() error

// Wraps a function withot error return type
func shutdownFunc(fn func()) shutdownHandler {
	return func() error {
		fn()
		return nil
	}
}

// Panic recovery and shutdown handler
func handleShutdown() {
	v := recover()
	if v != nil {
		logger.Error(v)
	}

	logger.Info("Shutting down")

	for _, handler := range shutdown {
		err := handler()
		if err != nil {
			logger.Error(err)
		}
	}

	if v != nil {
		panic(v)
	}
}
