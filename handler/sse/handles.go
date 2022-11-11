package sse

import (
	"net/http"
)

// Handle downstream
func Handle(w http.ResponseWriter, r *http.Request, h Streamer) {
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
	h.Pour(w, f)
}

// NewHandler returns a HandlerFunc that writes/loop the event to the ResponseWriter.
func NewHandler(h Streamer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r, h)
	}
}

// MsgHandle handle a signal message
type MsgHandle func(http.ResponseWriter, http.Flusher, any) error

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
