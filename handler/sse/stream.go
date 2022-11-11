package sse

import (
	"net/http"
	"time"
)

// Demo time tick SSE
func Demo(w http.ResponseWriter, r *http.Request) {
	tick := tickStream{}
	Handle(w, r, tick)
}

// Streamer write endless events to ResponseWriter.
type Streamer interface {
	Pour(http.ResponseWriter, http.Flusher)
}

type tickStream struct{}

// Pour ...
func (t tickStream) Pour(w http.ResponseWriter, f http.Flusher) {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for i := 0; i < 3; i++ {
		msg := NewMessage(i)
		err := msg.Send(w, f)
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
		err = NewMessage(jsonMsg).Send(w, f)
		if err != nil {
			return
		}
		<-timer.C
	}
	Shutdown(w, f)
}
