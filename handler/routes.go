package handler

import "net/http"

// Routes return http.Handler for http server
// must be call after add all custum path handler
func (g *RouterGroup) Routes(middlewares ...Middleware) http.Handler {
	g.r.router.NotFound = errResponse(http.StatusNotFound,
		"the requested resource could not be found",
	)
	g.r.router.MethodNotAllowed = MethodNotAllowed
	g.Get("/v1/healthcheck", HealthCheck)

	middlewares = append(middlewares, g.r.buildIns()...)
	h := http.Handler(g.r.router)
	mm := h.ServeHTTP
	for _, m := range middlewares {
		mm = m(mm)
	}
	return http.HandlerFunc(mm)
}
