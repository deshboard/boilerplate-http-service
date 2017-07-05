package main

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/deshboard/boilerplate-http-service/app"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/goph/serverz/named"
)

// newServer creates the main server instance for the service.
func newServer(appCtx *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	service := app.NewService()
	service.Logger = appCtx.logger
	service.ErrorHandler = appCtx.errorHandler

	handler := app.NewServiceHandler(service, appCtx.tracer)

	return &named.Server{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "http: ", 0),
		},
		ServerName: "http",
	}
}
