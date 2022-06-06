package handler

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

// router holds all paths relative funcs
type router struct {
	router *httprouter.Router
	config *Config
}

func (r *router) buildIns() []Middleware {
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
		r.router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
		ms = append(ms, r.metrics)
	}
	return ms
}
