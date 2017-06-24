package app

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

// Returns parameters from the request
// (decouples the service from the router implementation)
type paramFetcher func(r *http.Request) map[string]string

// Service contains the main controller logic
type Service struct {
	getParams paramFetcher

	logger log.Logger
}

// NewService creates a new service object
func NewService(logger log.Logger) *Service {
	return &Service{
		getParams: paramFetcher(mux.Vars),

		logger: logger,
	}
}

// Home represents the main route
func (s *Service) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
