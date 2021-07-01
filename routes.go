package gtea

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFountResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowResponse)

	router.HandlerFunc(
		http.MethodGet,
		"/v1/healthcheck",
		app.healthCheckHandler)

	if app.config.metrics {
		router.Handler(
			http.MethodGet,
			"/debug/vars",
			expvar.Handler())
	}
	rlMiddle := app.rateLimit(router)
	corsMiddle := app.enabledCORS(rlMiddle)
	recoverMiddle := app.recoverPanic(corsMiddle)
	return app.metrics(recoverMiddle)
}
