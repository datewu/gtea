package router

import (
	"expvar"
	"net/http"

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

func (r *bag) buildIns() []Middleware {
	ms := []Middleware{}
	// note the order is siginificant
	if r.config.Limiter.Enabled {
		ms = append(ms, r.rateLimit)
	}
	if r.config.CORS.TrustedOrigins != nil {
		ms = append(ms, r.enabledCORS)
	}
	ms = append(ms, recoverPanic)
	if r.config.Metrics {
		r.rt.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
		ms = append(ms, r.metrics)
	}
	return ms
}
