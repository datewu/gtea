package gtea

import (
	"fmt"
)

// Message chan feedback
type Message struct {
	Payload interface{}
	Err     error
}

// AddDaemonFn add defer func in app.shutdown you may add db.close, redis.close, etc
// not goroutine safe
func (app *App) AddDaemonFn(fn func()) {
	app.daemonWG.Add(1)
	app.daemonFns = append(app.daemonFns, fn)
}

// AddBGJob start a background job, goroutine safe
func (app *App) AddBGJob(name string, fn func(chan Message)) {
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.Err(fmt.Errorf("job %s, recoveed %s", name, r), nil)
		}
	}
	app.bgWG.Add(1)
	go func() {
		defer app.bgWG.Done()
		defer rcv()
		app.chansLock.Lock()
		c, ok := app.bgChans[name]
		if !ok {
			c = make(chan Message)
			app.bgChans[name] = c
		}
		app.chansLock.Unlock()
		fn(c)
	}()
}

// GetBGChan get backgroud job receive only feedback chan
// goroutine safe
func (app *App) GetBGChan(name string) <-chan Message {
	app.chansLock.Lock()
	defer app.chansLock.Unlock()
	return app.bgChans[name]
}

// RemoveBGChan remove job feedback chan
// goroutine safe
func (app *App) RemoveBGChan(name string) {
	app.chansLock.Lock()
	defer app.chansLock.Unlock()
	delete(app.bgChans, name)
}
