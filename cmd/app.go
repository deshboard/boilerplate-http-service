package main

import (
	"net/http"

	"github.com/deshboard/boilerplate-http-service/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	fxhttp "github.com/goph/fxt/http"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *app.Service {
	return app.NewService(
		app.Logger(params.Logger),
		app.ErrorHandler(params.ErrorHandler),
	)
}

// NewHandler constructs a new service handler instance.
func NewHandler(router *mux.Router, service *app.Service) http.Handler {
	return app.NewHandler(router, service)
}

// NewHTTPConfig creates a http config.
func NewHTTPConfig(config *Config) *fxhttp.Config {
	addr := config.HTTPAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := fxhttp.NewConfig(addr)

	return c
}
