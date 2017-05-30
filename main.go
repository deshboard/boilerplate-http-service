package main // import "github.com/deshboard/boilerplate-http-service"

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
	"github.com/deshboard/boilerplate-http-service/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
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

	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthService := healthz.NewHealthService(serviceChecker, status)
	healthHandler := http.NewServeMux()

	healthHandler.HandleFunc("/healthz", healthService.HealthStatus)
	healthHandler.HandleFunc("/readiness", healthService.ReadinessStatus)

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

	shutdownManager.RegisterAsFirst(healthServer.Close)
	go serverManager.ListenAndStartServer(healthServer, config.HealthAddr)(errChan)

	server := &serverz.NamedServer{
		Server: &http.Server{
			Handler:  newHandler(),
			ErrorLog: log.New(w, "http: ", 0),
		},
		Name: "http",
	}

	shutdownManager.RegisterAsFirst(server.Close)
	go serverManager.ListenAndStartServer(server, config.ServiceAddr)(errChan)

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

			go serverManager.StopServer(server, wg)(ctx)
			go serverManager.StopServer(healthServer, wg)(ctx)

			wg.Wait()

			// Cancel context if shutdown completed earlier
			cancel()

			// Break the loop, proceed with regular shutdown
			break MainLoop
		}
	}
}
