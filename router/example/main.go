package main

import (
	"fmt"
	"net/http"

	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/router"
)

func main() {
	conf := &router.Config{Debug: true}
	r := router.NewRouter(conf)

	r.Get("/ok", handler.HealthCheck)
	r.Get("/ok/good", handler.HealthCheck)
	r.Get("/ok/bye", handler.HealthCheck)
	r.Get("/ok/bye/lala", handler.HealthCheck)
	r.Get("/ok/bye/lala/a/b/c/d/e/f/z", handler.HealthCheck)
	r.Post("/ok", handler.HealthCheck)
	r.Put("/ok", handler.HealthCheck)
	r.Delete("/ok", handler.HealthCheck)
	r.Static("/abc", "../../")
	r.Static("/", "../")
	srv := &http.Server{
		Addr:    ":8090",
		Handler: r,
	}
	fmt.Println("start serve", srv.ListenAndServe())
}
