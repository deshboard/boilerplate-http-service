package app

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/gorilla/mux"
)

// Service contains the main controller logic
type Service struct {
	logger       log.Logger
	errorHandler emperror.Handler
}

// NewService creates a new service object
func NewService(logger log.Logger, errorHandler emperror.Handler) *Service {
	return &Service{logger, errorHandler}
}

// getParams returns parameters from the request.
// (decouples the service from the router implementation)
func (s *Service) getParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Home represents the main route
func (s *Service) Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
