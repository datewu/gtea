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

func (r *Router) buildIns() []handler.Middleware {
	ms := []handler.Middleware{}
	// note the order is siginificant
	// outer if executed first
	if r.conf.Limiter.Enabled {
		ms = append(ms, r.rateLimitMiddleware())
	}
	if r.conf.CORS.TrustedOrigins != nil {
		ms = append(ms, r.corsMiddleware())
	}
	ms = append(ms, handler.RecoverPanicMiddleware)
	if r.conf.Metrics {
		r.Handle(http.MethodGet, "/debug/vars", expvar.Handler())
		ms = append(ms, handler.MetricsMiddleware)
	}
	return ms
}
