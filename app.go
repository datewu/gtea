package gtea

import (
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/datewu/jsonlog"
)

// App is the main application object.
type App struct {
	Logger  *jsonlog.Logger
	config  *Config
	wg      sync.WaitGroup
	exitFns []func()
}

func NewApp(cfg *Config) *App {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	if cfg.Metrics {
		expvar.Publish("goroutines", expvar.Func(func() interface{} {
			return runtime.NumGoroutine()
		}))
		expvar.Publish("timestamp", expvar.Func(func() interface{} {
			return time.Now().Unix()
		}))
	}
	app := &App{
		config: cfg,
		Logger: logger,
	}
	return app
}
