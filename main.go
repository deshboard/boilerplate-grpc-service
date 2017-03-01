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
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func main() {
	defer shutdown.Handle()

	flag.Parse()

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
	}).Printf("Starting %s service", app.FriendlyServiceName)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdown.Register(w.Close)

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
		shutdown.RegisterAsFirst(debugServer.Close)

		go serverManager.ListenAndStartServer(debugServer, config.DebugAddr)(errChan)
	}

	grpcServer := grpc.NewServer()
	boilerplate.RegisterBoilerplateServer(grpcServer, app.NewService())

	healthHandler, status := newHealthServiceHandler()
	healthServer := &serverz.NamedServer{
		Server: &http.Server{
			Handler:  healthHandler,
			ErrorLog: log.New(w, "health: ", 0),
		},
		Name: "health",
	}
	shutdown.RegisterAsFirst(healthServer.Close, serverz.ShutdownFunc(grpcServer.Stop))

	go serverManager.ListenAndStartServer(healthServer, config.HealthAddr)(errChan)

	go func() {
		lis, err := net.Listen("tcp", config.ServiceAddr)
		if err != nil {
			errChan <- err
			return
		}

		logger.WithField("addr", lis.Addr()).WithField("server", "grpc").Info("Server started")
		errChan <- grpcServer.Serve(lis)
	}()

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

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			wg := &sync.WaitGroup{}

			if config.Debug {
				go serverManager.StopServer(debugServer, wg)(ctx)
			}
			go serverManager.StopServer(healthServer, wg)(ctx)

			wg.Add(1)

			go func() {
				// TODO: implement timeout
				grpcServer.GracefulStop()

				wg.Done()
			}()

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}

	close(errChan)
	close(signalChan)
}
