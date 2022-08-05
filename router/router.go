package router

import (
	"expvar"
	"net/http"

	"github.com/datewu/gtea/handler"
	"github.com/julienschmidt/httprouter"
)

// Config is the configuration for the router
type Config struct {
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	CORS struct {
		TrustedOrigins []string
	}
	Metrics bool
}

// DefaultConf return the default config
func DefaultConf() *Config {
	cnf := &Config{Metrics: true}
	cnf.Limiter.Enabled = true
	cnf.Limiter.Rps = 200
	cnf.Limiter.Burst = 10
	cnf.CORS.TrustedOrigins = nil
	return cnf
}

// bag holds all paths relative funcs
type bag struct {
	rt     *httprouter.Router
	config *Config
}

func (r *bag) buildIns() []handler.Middleware {
	ms := []handler.Middleware{}
	// note the order is siginificant
	if r.config.Limiter.Enabled {
		ms = append(ms, r.rateLimitMiddleware())
	}
	if r.config.CORS.TrustedOrigins != nil {
		ms = append(ms, r.corsMiddleware())
	}
	ms = append(ms, handler.RecoverPanicMiddleware)
	if r.config.Metrics {
		r.rt.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
		ms = append(ms, handler.MetricsMiddleware)
	}
	return ms
}
