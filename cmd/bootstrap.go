package main

import (
	"net/http"

	"github.com/deshboard/boilerplate-http-service/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	fxhttp "github.com/goph/fxt/http"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService constructs a new service instance.
func NewService(params ServiceParams) *app.Service {
	return app.NewService(
		app.Logger(params.Logger),
		app.ErrorHandler(params.ErrorHandler),
	)
}

// NewService constructs a new service handler instance.
func NewServiceHandler(service *app.Service, tracer opentracing.Tracer) http.Handler {
	router := app.NewRouter(tracer)

	router.HandleFunc("/", service.Home).Name("index").Methods("GET")

	return router
}

// NewHTTPConfig creates a http config.
func NewHTTPConfig(config *Config) *fxhttp.Config {
	c := fxhttp.NewConfig(config.HTTPAddr)

	return c
}
