package router

// Config is the configuration for the router
type Config struct {
	Debug   bool
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
	cnf := &Config{Metrics: true, Debug: true}
	cnf.Limiter.Enabled = true
	cnf.Limiter.Rps = 200
	cnf.Limiter.Burst = 10
	cnf.CORS.TrustedOrigins = nil
	return cnf
}
