package gtea

import (
	"database/sql"
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/datewu/jsonlog"
)

type application struct {
	Logger *jsonlog.Logger
	config *Config
	wg     sync.WaitGroup
	db     *sql.DB
}

func NewApp(cfg *Config) *application {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
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
		Logger: logger,
	}
	return app
}

func (app *application) SetDB(db *sql.DB) {
	app.db = db
}

func (app *application) shutdown() {
	if app.db != nil {
		err := app.db.Close()
		if err != nil {
			app.Logger.PrintErr(err, nil)
		}
	}
}
