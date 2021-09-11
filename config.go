package gtea

// Config is the configuration for the application
type Config struct {
	Port    int
	Env     string
	Metrics bool
}

// DefaultConfig is the default configuration for the application
func DefaultConfig() *Config {
	return &Config{
		Port:    8080,
		Env:     "development",
		Metrics: false,
	}
}
