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
func (app *App) Serve(routes http.Handler) error {
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
		app.Logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}
		app.Logger.PrintInfo("completing background tasks", nil)
		app.wg.Wait()
		shutdownErr <- nil
	}
	go bgSignal()
	app.Logger.PrintInfo("starting server", map[string]string{
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
	app.Logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
