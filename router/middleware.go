package router

import (
	"expvar"
	"net/http"

	"github.com/datewu/gtea/handler"
)

func (r *Router) rateLimitMiddleware() handler.Middleware {
	return handler.RateLimitMiddleware(r.conf.Limiter.Rps, r.conf.Limiter.Burst)
}

func (r *Router) corsMiddleware() handler.Middleware {
	return handler.CORSMiddleware(r.conf.CORS.TrustedOrigins)
}

func (r *Router) aggBuildInMiddlewares() {
	if r.conf.Limiter.Enabled {
		r.middleware = handler.Insert(r.middleware, r.rateLimitMiddleware())
	}
	if r.conf.CORS.TrustedOrigins != nil {
		r.middleware = handler.Insert(r.middleware, r.corsMiddleware())
	}
	if r.conf.Metrics {
		r.Handle(http.MethodGet, "/debug/vars", expvar.Handler())
		r.middleware = handler.Insert(r.middleware, handler.MetricsMiddleware)
	}
	r.middleware = handler.Append(r.middleware, handler.RecoverPanicMiddleware)
}
