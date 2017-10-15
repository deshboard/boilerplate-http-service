package app

import (
	"net/http"

	fxhttp "github.com/goph/fxt/http"
	"github.com/opentracing/opentracing-go"
)

// NewServiceHandler returns a new service handler
// IMPORTANT: all routes SHOULD have a name
func NewServiceHandler(service *Service, tracer opentracing.Tracer) http.Handler {
	router := fxhttp.NewRouter(tracer)

	router.HandleFunc("/", service.Home).Name("index").Methods("GET")

	return router
}
