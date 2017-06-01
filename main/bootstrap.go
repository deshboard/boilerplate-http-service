package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-http-service/app"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func bootstrap() serverz.Server {
	serviceChecker := healthz.NewTCPChecker(config.ServiceAddr, healthz.WithTCPTimeout(2*time.Second))
	checkerCollector.RegisterChecker(healthz.LivenessCheck, serviceChecker)

	service := app.NewService(logger)

	handler := app.NewServiceHandler(service, tracer)

	w := logger.Logger.WriterLevel(logrus.ErrorLevel)
	shutdownManager.Register(w.Close)

	return &serverz.NamedServer{
		Server: &http.Server{
			Handler:  handler,
			ErrorLog: log.New(w, "http: ", 0),
		},
		Name: "http",
	}
}
