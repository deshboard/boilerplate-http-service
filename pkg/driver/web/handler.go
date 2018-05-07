package web

import "github.com/gorilla/mux"

// RegisterHandlers registers the service in the HTTP handler mux.
// IMPORTANT: all routes SHOULD have a name
func RegisterHandlers(router *mux.Router, service *Service) {
	router.HandleFunc("/", service.Home).Name("index").Methods("GET")
}
