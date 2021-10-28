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
	Logger   *jsonlog.Logger
	metaData map[string]string
	config   *Config
	wg       sync.WaitGroup
	exitFns  []func()
}

// AddMetaData adds meta data to the application by key.
// skipped if the key already exists.
func (a *App) AddMetaData(key, value string) {
	if _, ok := a.metaData[key]; ok {
		return
	}
	a.metaData[key] = value
}

// GetMetaData returns the meta data of the application by key.
func (a *App) GetMetaData(key string) string {
	return a.metaData[key]
}

// NewApp creates a new application object.
func NewApp(cfg *Config) *App {
	logger := jsonlog.New(os.Stdout, cfg.LogLevel)
	if cfg.Metrics {
		expvar.Publish("goroutines", expvar.Func(func() interface{} {
			return runtime.NumGoroutine()
		}))
		expvar.Publish("timestamp", expvar.Func(func() interface{} {
			return time.Now().Unix()
		}))
	}
	app := &App{
		metaData: make(map[string]string),
		config:   cfg,
		Logger:   logger,
	}
	return app
}

// DefaultApp is the default application object.
func DefaultApp() *App {
	cfg := DefaultConfig()
	return NewApp(cfg)
}
