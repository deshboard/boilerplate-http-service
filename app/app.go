package app

import (
	webfx "github.com/deshboard/boilerplate-http-service/pkg/driver/fx"
	"github.com/goph/fxt/app/http"
	"github.com/goph/fxt/http/gorilla/mux"
	"github.com/goph/fxt/http/gorilla/mux/opentracing"
	"github.com/goph/fxt/tracing/opentracing"
	"go.uber.org/fx"
)

// Module is the collection of all modules of the application.
var Module = fx.Options(
	fxhttpapp.Module,

	// Configuration
	fx.Provide(
		NewLoggerConfig,
		NewDebugConfig,
	),

	// HTTP server
	fx.Provide(
		fxmux.Module,
		NewHTTPConfig,
	),

	fx.Provide(fxopentracing.NewTracer),

	webfx.Module,

	// Make sure to register this invoke function as the last,
	// so tracer is injected into all routes.
	fx.Invoke(otmux.InjectTracer),
)

type Runner = fxhttpapp.Runner
