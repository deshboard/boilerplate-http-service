package app

import fxhttp "github.com/goph/fxt/http"

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
