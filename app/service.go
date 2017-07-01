package app

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/gorilla/mux"
)

// Service contains the main controller logic.
type Service struct {
	Logger       log.Logger
	ErrorHandler emperror.Handler
}

// NewService creates a new service object.
func NewService() *Service {
	return &Service{
		Logger:       log.NewNopLogger(),
		ErrorHandler: emperror.NewNullHandler(),
	}
}

// getParams returns parameters from the request.
// (decouples the service from the router implementation)
func (s *Service) getParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Home represents the main route.
func (s *Service) Home(w http.ResponseWriter, r *http.Request) {
	level.Info(s.Logger).Log("msg", "Home page loaded")

	w.WriteHeader(200)
}
