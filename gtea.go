package gtea

import (
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/datewu/jsonlog"
)

// Config is the configuration for the application
type Config struct {
	Port    int
	Env     string
	Metrics bool
}

type application struct {
	config *Config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func NewApp(cfg *Config) *application {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	logger.PrintInfo("database connection pool established", nil)

	if cfg.Metrics {
		expvar.Publish("goroutines", expvar.Func(func() interface{} {
			return runtime.NumGoroutine()
		}))
		expvar.Publish("timestamp", expvar.Func(func() interface{} {
			return time.Now().Unix()
		}))
	}
	app := &application{
		config: cfg,
		logger: logger,
	}
	return app
}
