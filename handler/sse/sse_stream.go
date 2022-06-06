package api

import (
	"bytes"
	"encoding/json"
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
	w.WriteHeader(http.StatusOK)

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	h(w, f)
}

// Event is a single event in an SSE stream.
// use json marshal data
type Event struct {
	id, name *string
	data     any
	retry    *string
}

// Bytes marshal the event as a byte slice.
// special encode '\n' to '\r\n'
func (s Event) Bytes() []byte {
	if s.data == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	if s.id != nil {
		buf.WriteString(`id: `)
		buf.WriteString(*s.id + "\n")
	}
	if s.name != nil {
		buf.WriteString(`event: `)
		buf.WriteString(*s.name + "\n")
	}
	if s.retry != nil {
		buf.WriteString(`retry: `)
		buf.WriteString(*s.retry + "\n")
	}
	if s.data != nil {
		buf.WriteString(`data: `)
		bs, err := json.Marshal(s.data)
		if err != nil {
			return nil
		}
		for _, b := range bs {
			if b == '\n' {
				buf.WriteString("\r\n")
			} else {
				buf.WriteByte(b)
			}
		}
	}
	buf.WriteString("\n\n")
	return buf.Bytes()
}

// NewEvent creates a new SSE event.
func NewEvent(name string, data any) Event {
	return Event{
		name: &name,
		data: data,
	}
}

// Message is a single message in an SSE stream.
func NewMessage(data any) Event {
	return Event{
		data: data,
	}
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
