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

func (app *App) httpServer(ctx context.Context, routes http.Handler) *http.Server {
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.Port),
		Handler:  routes,
		ErrorLog: log.New(app.Logger, "", 0),
	}
	srv.IdleTimeout = time.Minute
	srv.ReadTimeout = 10 * time.Second
	if !app.config.NoWirteTimeout {
		srv.WriteTimeout = 30 * time.Second
	}
	return srv
}

// Serve start https server with cert private PEM encode file
func (app *App) ServeTLS(ctx context.Context, routes http.Handler, cert, pKey string) error {
	srv := app.httpServer(ctx, routes)
	s := func() error {
		return srv.ListenAndServeTLS(cert, pKey)
	}
	return app.serve(ctx, srv, s)
}

// Serve start http server
func (app *App) serve(ctx context.Context, srv *http.Server, s func() error) error {
	go handleOSsignal(ctx, app, srv)
	app.Logger.Info("starting server", map[string]any{
		"env":  app.config.Env,
		"addr": srv.Addr,
	})
	err := s()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-app.shutdownStream
	if err != nil {
		return err
	}
	app.Logger.Info("stopped server", map[string]any{
		"addr": srv.Addr,
	})
	app.Shutdown()
	return nil
}

// Serve start http server
func (app *App) Serve(ctx context.Context, routes http.Handler) error {
	srv := app.httpServer(ctx, routes)
	return app.serve(ctx, srv, srv.ListenAndServe)
}

// Shutdown call clearFns one by one
func (app *App) Shutdown() {
	for _, fn := range app.clearFns {
		app.clearWG.Add(1)
		go fn()
	}
	app.clearWG.Wait()
	app.clearFns = nil
}

func handleOSsignal(ctx context.Context, app *App, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	app.Logger.Info("caught signal", map[string]any{
		"signal": s.String(),
	})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		app.shutdownStream <- err
	}
	app.Logger.Info("going cancel all background jobs")
	for name, j := range app.bgJobs {
		j.Cancle()
		app.Logger.Info("canceled bg job:", map[string]any{"name": name})
	}
	app.bgWG.Wait()
	app.Logger.Info("all bgWG wait done")
	app.shutdownStream <- nil
	close(app.shutdownStream)
}
