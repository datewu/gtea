package gtea

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/gtea/router"
)

func TestIntegrate(t *testing.T) {
	port := 32808
	cfg := Config{
		Port:     port,
		Env:      "development",
		Metrics:  true,
		LogLevel: jsonlog.LevelDebug,
	}
	ctx := context.Background()
	app := NewApp(ctx, &cfg)
	h := router.DefaultRoutesGroup()
	h.Get("/", handler.HealthCheck)
	h.Get("/api/v1", handler.HealthCheck)
	go app.Serve(ctx, h.Handler())
	fn := func(path string) {
		url := fmt.Sprintf("http://localhost:%d%s", port, path)
		res, err := http.Get(url)
		if err != nil {
			t.Error(err)
			return
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}
		got := string(data)
		if got != `{"status":"available"}` {
			t.Errorf("expect %q, got %q", `{"status":"available"}`, got)
		}
	}
	go fn("/")
	go fn("/api/v1")
	go fn("/api/v1")
	go fn("/api/v1")
	time.Sleep(2 * time.Second)
	app.Shutdown()
}
