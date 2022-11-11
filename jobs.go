package gtea

import (
	"context"
	"fmt"
)

// Message chan feedback
type Message struct {
	Payload interface{}
	Err     error
}

// JobParam ...
type JobParam struct {
	Chan   chan Message
	Cancle context.CancelFunc
}

// AddClearFn add defer func in app.shutdown you may add db.close, redis.close, etc
// not goroutine safe
func (app *App) AddClearFn(fn func()) {
	app.clearWG.Add(1)
	f := func() {
		fn()
		app.clearWG.Done()
	}
	app.clearFns = append(app.clearFns, f)
}

// AddBGJob start a background job, goroutine safe
func (app *App) AddBGJob(name string, fn func(context.Context, chan<- Message)) error {
	app.bgLock.Lock()
	_, ok := app.bgJobs[name]
	app.bgLock.Unlock()
	if ok {
		return fmt.Errorf("cannot overrider job %s", name)
	}
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.Err(fmt.Errorf("job %s, recoveed %s", name, r), nil)
		}
	}
	app.bgWG.Add(1)
	go func() {
		defer app.bgWG.Done()
		defer rcv()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		param := &JobParam{
			Chan:   make(chan Message),
			Cancle: cancel,
		}
		app.bgLock.Lock()
		app.bgJobs[name] = param
		app.bgLock.Unlock()
		fn(ctx, param.Chan)
	}()
	return nil
}

// GetBGJobParam get backgroud job param
// goroutine safe
func (app *App) GetBGJobParam(name string) *JobParam {
	app.bgLock.Lock()
	defer app.bgLock.Unlock()
	return app.bgJobs[name]
}

// GetBGChan get backgroud job receive only feedback chan
// goroutine safe
func (app *App) GetBGChan(name string) <-chan Message {
	app.bgLock.Lock()
	job := app.bgJobs[name]
	app.bgLock.Unlock()
	if job == nil {
		return nil
	}
	return job.Chan
}

// RemoveBGChan remove job feedback chan
// goroutine safe
func (app *App) RemoveBGChan(name string) {
	app.bgLock.Lock()
	defer app.bgLock.Unlock()
	delete(app.bgJobs, name)
}
