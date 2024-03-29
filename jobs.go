package gtea

import (
	"context"
	"errors"
	"fmt"
)

// Message chan feedback
type Message struct {
	Payload any
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
	f := func() {
		defer app.clearWG.Done()
		fn()
	}
	app.clearFns = append(app.clearFns, f)
}

// FireAnonymousJob start a background job, goroutine safe
func (app *App) FireAnonymousJob(fn func(context.Context)) {
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.Err(errors.New("an anonymous job, recovered"))
		}
	}
	app.bgWG.Add(1)
	go func() {
		defer app.bgWG.Done()
		ctx, cancel := context.WithCancel(app.ctx)
		defer cancel()
		defer rcv()
		fn(ctx)
	}()
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
		defer app.removeBGJob(name)
		defer app.bgWG.Done()
		ctx, cancel := context.WithCancel(app.ctx)
		defer cancel()
		param := &JobParam{Chan: make(chan Message), Cancle: cancel}
		app.bgLock.Lock()
		app.bgJobs[name] = param
		app.bgLock.Unlock()
		defer rcv()
		defer close(param.Chan)
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

// removeBGJob remove background job
// goroutine safe
func (app *App) removeBGJob(name string) {
	app.bgLock.Lock()
	defer app.bgLock.Unlock()
	delete(app.bgJobs, name)
}
