package main

import (
	"net/http"

	"github.com/deshboard/boilerplate-http-service/app"
)

func newHandler() http.Handler {
	service := app.NewService(logger)

	return app.NewServiceHandler(service, tracer)
}
