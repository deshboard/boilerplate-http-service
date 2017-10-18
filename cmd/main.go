package main

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/fxt"
	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/errors"
	"github.com/goph/fxt/http"
	"github.com/goph/fxt/http/gorilla"
	fxlog "github.com/goph/fxt/log"
	"github.com/goph/fxt/tracing"
	"github.com/goph/healthz"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func main() {
	status := healthz.NewStatusChecker(healthz.Healthy)
	var ext struct {
		Config       *Config
		Closer       fxt.Closer
		Logger       log.Logger
		ErrorHandler emperror.Handler

		DebugErr debug.Err
		HTTPErr  http.Err
	}

	app := fx.New(
		fx.NopLogger,
		fxt.Bootstrap,
		fx.Provide(
			NewConfig,

			// Log and error handling
			NewLoggerConfig,
			fxlog.NewLogger,
			errors.NewHandler,

			// Debug server
			NewDebugConfig,
			debug.NewServer,
			debug.NewHealthCollector,
		),
		fx.Invoke(func(collector healthz.Collector) {
			collector.RegisterChecker(healthz.ReadinessCheck, status)
		}),
		fx.Extract(&ext),

		fx.Provide(
			// HTTP server
			mux.NewRouter,
			NewService,
			NewHandler,
			NewHTTPConfig,
			http.NewServer,

			tracing.NewTracer,
		),
		fx.Invoke(gorilla.InjectTracer),
	)

	err := app.Err()
	if err != nil {
		panic(err)
	}

	// Close resources when the application stops running
	defer ext.Closer.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(ext.ErrorHandler)

	level.Info(ext.Logger).Log(
		"msg", fmt.Sprintf("starting %s", FriendlyServiceName),
		"version", Version,
		"commit_hash", CommitHash,
		"build_date", BuildDate,
	)

	err = app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	select {
	case sig := <-app.Done():
		level.Info(ext.Logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-ext.DebugErr:
		if err != nil {
			err = emperror.WithStack(emperror.WithMessage(err, "debug server crashed"))
			ext.ErrorHandler.Handle(err)
		}

	case err := <-ext.HTTPErr:
		if err != nil {
			err = emperror.WithStack(emperror.WithMessage(err, "http server crashed"))
			ext.ErrorHandler.Handle(err)
		}
	}

	status.SetStatus(healthz.Unhealthy)

	ctx, cancel := context.WithTimeout(context.Background(), ext.Config.ShutdownTimeout)
	defer cancel()

	err = app.Stop(ctx)
	emperror.HandleIfErr(ext.ErrorHandler, err)
}
