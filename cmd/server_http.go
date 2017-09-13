package main

import (
	stdlog "log"
	"net/http"
	"time"

	. "github.com/deshboard/boilerplate-http-service/app"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
)

// newHTTPServer creates the main server instance for the service.
func newHTTPServer(app *application) serverz.Server {
	serviceChecker := healthz.NewTCPChecker(app.config.HTTPAddr, healthz.WithTCPTimeout(2*time.Second))
	app.healthCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	service := NewService(
		Logger(app.logger),
		ErrorHandler(app.errorHandler),
	)

	handler := NewServiceHandler(service, app.tracer)

	return &serverz.AppServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: stdlog.New(log.NewStdlibAdapter(level.Error(app.logger)), "http: ", 0),
		},
		Name:   "http",
		Addr:   serverz.NewAddr("tcp", app.config.HTTPAddr),
		Logger: app.logger,
	}
}