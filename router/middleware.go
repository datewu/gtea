package router

import (
	"github.com/datewu/gtea/handler"
)

func (ro *bag) rateLimitMiddleware() handler.Middleware {
	return handler.RateLimit(ro.config.Limiter.Rps, ro.config.Limiter.Burst)
}

func (ro *bag) corsMiddleware() handler.Middleware {
	return handler.CORS(ro.config.CORS.TrustedOrigins)
}
