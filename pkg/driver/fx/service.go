package fx

import (
	"github.com/deshboard/boilerplate-http-service/pkg/driver/web"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *web.Service {
	return web.NewService(
		web.Logger(params.Logger),
		web.ErrorHandler(params.ErrorHandler),
	)
}
