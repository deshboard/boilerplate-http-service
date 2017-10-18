package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewHandler returns a new HTTP handler.
// IMPORTANT: all routes SHOULD have a name
func NewHandler(router *mux.Router, service *Service) http.Handler {
	router.HandleFunc("/", service.Home).Name("index").Methods("GET")

	return router
}
