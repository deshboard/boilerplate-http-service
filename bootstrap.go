package main

import (
	"net/http"

	"github.com/deshboard/boilerplate-http-service/app"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sagikazarmark/healthz"
)

// Returns a new service handler
// IMPORTANT: all routes SHOULD have a name
func newServiceHandler(service *app.Service, tracer opentracing.Tracer) http.Handler {
	router := app.NewRouter(tracer)

	router.HandleFunc("/", service.Home).Name("index").Methods("GET")

	return router
}

// Creates the health service handler and the status checker
func newHealthServiceHandler() (http.Handler, *healthz.StatusChecker) {
	status := healthz.NewStatusChecker(healthz.Healthy)
	healthMux := healthz.NewHealthServiceHandler(healthz.NewCheckers(), status)

	return healthMux, status
}
