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
	Logger         *jsonlog.Logger
	shutdownStream chan error
	clearWG        sync.WaitGroup
	clearFns       []func()
	bgLock         *sync.Mutex
	bgWG           sync.WaitGroup
	bgJobs         map[string]*JobParam
}

// NewApp creates a new application object.
func NewApp(ctx context.Context, cfg *Config) *App {
	logger := jsonlog.New(os.Stdout, cfg.LogLevel)
	if cfg.Metrics {
		expvar.Publish("goroutines", expvar.Func(func() any {
			return runtime.NumGoroutine()
		}))
		expvar.Publish("timestamp", expvar.Func(func() any {
			return time.Now().Unix()
		}))
	}
	app := &App{
		ctx:            ctx,
		config:         cfg,
		shutdownStream: make(chan error),
		bgLock:         &sync.Mutex{},
		bgJobs:         make(map[string]*JobParam),
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
	Port           int
	Env            string
	Metrics        bool
	LogLevel       jsonlog.Level
	NoWirteTimeout bool
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

// AddMetaData adds meta data to the application by key.
func (a *App) AddMetaData(key string, value string) {
	ctx := context.WithValue(a.ctx, contextKey(key), value)
	a.ctx = ctx
}

// GetMetaData returns the meta data of the application by key.
func (a *App) GetMetaData(key string) string {
	v := a.ctx.Value(contextKey(key))
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}
