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

	shutdownErr := make(chan error)
	bgSignal := func() {
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
			shutdownErr <- err
		}
		app.Logger.Info("completing background tasks", nil)
		app.wg.Wait()
		shutdownErr <- nil
	}
	go bgSignal()
	app.Logger.Info("starting server", map[string]string{
		"env":  app.config.Env,
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownErr
	if err != nil {
		return err
	}
	app.Logger.Info("stopped server", map[string]string{
		"addr": srv.Addr,
	})
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
