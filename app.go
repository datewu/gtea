package gtea

import (
	"context"
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/datewu/gtea/jsonlog"
)

type contextKey string

// App is the main application object.
type App struct {
	ctx            context.Context
	config         *Config
	wg             sync.WaitGroup
	Logger         *jsonlog.Logger
	shutdownStream chan error
	exitFns        []func()
}

// NewApp creates a new application object.
func NewApp(ctx context.Context, cfg *Config) *App {
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
		ctx:            ctx,
		config:         cfg,
		shutdownStream: make(chan error),
		Logger:         logger,
	}
	return app
}

// DefaultApp is the default application object.
func DefaultApp() *App {
	cfg := DefaultConfig()
	return NewApp(context.Background(), cfg)
}

// Config is the configuration for the application
type Config struct {
	Port     int
	Env      string
	Metrics  bool
	LogLevel jsonlog.Level
}

// DefaultConfig is the default configuration for the application
func DefaultConfig() *Config {
	return &Config{
		Port:     8080,
		Env:      "development",
		Metrics:  false,
		LogLevel: jsonlog.LevelInfo,
	}
}
