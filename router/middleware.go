package router

import (
	"github.com/datewu/gtea/handler"
)

func (ro *bag) rateLimitMiddleware() handler.Middleware {
	return handler.RateLimitMiddleware(ro.config.Limiter.Rps, ro.config.Limiter.Burst)
}

func (ro *bag) corsMiddleware() handler.Middleware {
	return handler.CORSMiddleware(ro.config.CORS.TrustedOrigins)
}
