package app

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// NewServiceHandler returns a new service handler
// IMPORTANT: all routes SHOULD have a name
func NewServiceHandler(service *Service, tracer opentracing.Tracer) http.Handler {
	router := NewRouter(tracer)

	router.HandleFunc("/", service.Home).Name("index").Methods("GET")

	return router
}
