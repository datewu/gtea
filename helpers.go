package gtea

import (
	"fmt"
)

// Background start bg job
func (app *App) Background(fn func()) {
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.Err(fmt.Errorf("%s", r), nil)
		}
	}
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer rcv()
		fn()
	}()
}

// AddExitFn add defer func in app.shutdown
// you may add db.close, redis.close, etc
func (app *App) AddExitFn(fn func()) {
	app.exitFns = append(app.exitFns, fn)
}

// Shutdown shutdown app
func (app *App) Shutdown() {
	for _, fn := range app.exitFns {
		fn()
	}
}
