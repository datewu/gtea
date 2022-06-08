package sse

import (
	"net/http"
	"time"
)

// NewHandler returns a HandlerFunc that writes/loop the event to the ResponseWriter.
func NewHandler(h Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hanlder(w, r, h)
	}
}

// Hook is a SSE hook.
type Hook func(http.ResponseWriter, http.Flusher)

// Demo for SSE
func Demo(w http.ResponseWriter, r *http.Request) {
	hanlder(w, r, demoLoop)
}

func hanlder(w http.ResponseWriter, r *http.Request, h Hook) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	h(w, f)
}

func demoLoop(w http.ResponseWriter, f http.Flusher) {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for i := 0; i < 3; i++ {
		msg := NewMessage(i)
		w.Write(msg.Bytes())
		jsonMsg := struct {
			ID   int
			Time time.Time
		}{
			ID:   i,
			Time: time.Now(),
		}
		w.Write(NewMessage(jsonMsg).Bytes())
		f.Flush()
		<-timer.C
	}
	shutMsg := NewEvent("shutdown", "bye")
	w.Write(shutMsg.Bytes())
	f.Flush()
}
