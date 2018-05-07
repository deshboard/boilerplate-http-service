package fx

import (
	"github.com/deshboard/boilerplate-http-service/pkg/driver/web"
	"go.uber.org/fx"
)

// Module is an fx compatible module.
var Module = fx.Options(
	fx.Provide(NewService),
	fx.Invoke(web.RegisterHandlers),
)
