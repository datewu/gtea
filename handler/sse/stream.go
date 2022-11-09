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

// MsgHandle handle a signal message
type MsgHandle func(http.ResponseWriter, http.Flusher, any) error

// Demo time tick SSE
func Demo(w http.ResponseWriter, r *http.Request) {
	Handle(w, r, demoLoop)
}

// SendStringMsg sugar
func SendStringMsg(w http.ResponseWriter, f http.Flusher, msg string) error {
	eMsg := NewMessage(msg)
	_, err := w.Write(eMsg.Bytes())
	if err != nil {
		return err
	}
	f.Flush()
	return nil
}

// SendAnyMsg sugar
func SendAnyMsg(w http.ResponseWriter, f http.Flusher, msg interface{}) error {
	eMsg := NewMessage(msg)
	_, err := w.Write(eMsg.Bytes())
	if err != nil {
		return err
	}
	f.Flush()
	return nil
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

// Shutdown send shutdown event to client
func Shutdown(w http.ResponseWriter, f http.Flusher) error {
	shutMsg := NewEvent("shutdown", "bye")
	_, err := w.Write(shutMsg.Bytes())
	if err != nil {
		return err
	}
	f.Flush()
	return nil
}

func demoLoop(w http.ResponseWriter, f http.Flusher) {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for i := 0; i < 3; i++ {
		msg := NewMessage(i)
		_, err := w.Write(msg.Bytes())
		if err != nil {
			return
		}
		jsonMsg := struct {
			ID   int
			Time time.Time
		}{
			ID:   i,
			Time: time.Now(),
		}
		SendAnyMsg(w, f, jsonMsg)
		<-timer.C
	}
	Shutdown(w, f)
}
