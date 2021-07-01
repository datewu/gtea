package gtea

import (
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/datewu/jsonlog"
)

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
	metrics bool
}

type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func newTodo(cfg config) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	logger.PrintInfo("database connection pool established", nil)

	if cfg.metrics {
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

	err := app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
