package gtea

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleOSsignal(ctx context.Context, app *App, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	app.Logger.Info("caught signal", map[string]string{
		"signal": s.String(),
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		app.shutdownStream <- err
	}
	app.Logger.Info("completing background tasks", nil)
	app.wg.Wait()
	app.shutdownStream <- nil
	close(app.shutdownStream)
}

// Serve start http server
func (app *App) Serve(ctx context.Context, routes http.Handler) error {
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.Port),
		Handler:  routes,
		ErrorLog: log.New(app.Logger, "", 0),
	}
	srv.IdleTimeout = time.Minute
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 30 * time.Second

	go handleOSsignal(ctx, app, srv)
	app.Logger.Info("starting server", map[string]string{
		"env":  app.config.Env,
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-app.shutdownStream
	if err != nil {
		return err
	}
	app.Logger.Info("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	app.shutdown()
	return nil
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

// AddBGJob start a background job, goroutine safe
func (app *App) AddBGJob(name string, fn func(chan string)) {
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.Err(fmt.Errorf("%s", r), nil)
		}
	}
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer rcv()
		app.chansLock.Lock()
		c, ok := app.bgChans[name]
		if !ok {
			c = make(chan string)
			app.bgChans[name] = c
		}
		app.chansLock.Unlock()
		fn(c)
	}()
}

// GetBGChan get backgroud job feedback chan, goroutine safe
func (app *App) GetBGChan(name string) chan string {
	app.chansLock.Lock()
	defer app.chansLock.Unlock()
	return app.bgChans[name]
}

// RemoveBGChan remove job feedback chan, goroutine safe
func (app *App) RemoveBGChan(name string) {
	app.chansLock.Lock()
	defer app.chansLock.Unlock()
	delete(app.bgChans, name)
}

// AddExitFn add defer func in app.shutdown you may add db.close, redis.close, etc
// not goroutine safe
func (app *App) AddExitFn(fn func()) {
	app.exitFns = append(app.exitFns, fn)
}

func (app *App) shutdown() {
	for _, fn := range app.exitFns {
		fn()
	}
}
