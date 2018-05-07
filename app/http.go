package app

import (
	"net/http"

	fxhttp "github.com/goph/fxt/http"
	"github.com/gorilla/mux"
)

// NewHandler constructs a new HTTP handler instance.
func NewHandler(router *mux.Router) http.Handler {
	return router
}

// NewHTTPConfig creates a http config.
func NewHTTPConfig(config Config) *fxhttp.Config {
	addr := config.HTTPAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := fxhttp.NewConfig(addr)

	return c
}
