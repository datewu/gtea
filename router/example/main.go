package main

import (
	"net/http"

	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/router"
)

func main() {
	conf := &router.Config{}
	r := router.NewRouter(conf)

	r.Get("/ok", handler.HealthCheck)
	r.Get("/ok/good", handler.HealthCheck)
	r.Get("/ok/bye", handler.HealthCheck)
	r.Get("/ok/bye/lala", handler.HealthCheck)
	r.Get("/ok/bye/lala/a/b/c/d/e/f/z", handler.HealthCheck)
	r.Post("/ok", handler.HealthCheck)
	r.Put("/ok", handler.HealthCheck)
	r.Delete("/ok", handler.HealthCheck)
	r.Static("/abc", "../")
	r.Debug()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()
}
