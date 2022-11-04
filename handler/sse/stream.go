package sse

import (
	"net/http"
	"time"
)

// NewHandler returns a HandlerFunc that writes/loop the event to the ResponseWriter.
func NewHandler(h Downstream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r, h)
	}
}

// Downstream take over the responsibility to write the event to the ResponseWriter.
type Downstream func(http.ResponseWriter, http.Flusher)

// Demo time tick SSE
func Demo(w http.ResponseWriter, r *http.Request) {
	Handle(w, r, demoLoop)
}

// SendStringMsg sugar
func SendStringMsg(w http.ResponseWriter, f http.Flusher, msg string) {
	eMsg := NewMessage(msg)
	w.Write(eMsg.Bytes())
	f.Flush()
}

// Handle downstream
func Handle(w http.ResponseWriter, r *http.Request, h Downstream) {
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

// Shutdown send shutdown event
func Shutdown(w http.ResponseWriter, f http.Flusher) {
	shutMsg := NewEvent("shutdown", "bye")
	w.Write(shutMsg.Bytes())
	f.Flush()
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
	Shutdown(w, f)
}
