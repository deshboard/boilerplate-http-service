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
)

// newHTTPServer creates the main server instance for the service.
func newHTTPServer(appCtx *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(appCtx.config.HTTPAddr, healthz.WithTCPTimeout(2*time.Second))
	appCtx.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	service := app.NewService(
		app.Logger(appCtx.logger),
		app.ErrorHandler(appCtx.errorHandler),
	)

	handler := app.NewServiceHandler(service, appCtx.tracer)

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(appCtx.logger)), "http: ", 0),
		},
		Name:   "http",
		Addr:   serverz.NewAddr("tcp", appCtx.config.HTTPAddr),
		Logger: appCtx.logger,
	}
}
