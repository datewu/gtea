package sse

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Event a single event in SSE stream using json marshal encoding
type Event struct {
	id, name *string
	data     any
	retry    *string
}

// Bytes marshal the event as a byte slice. Escape '\n' to '\r\n'
func (e Event) Bytes() []byte {
	if e.data == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	if e.id != nil {
		buf.WriteString("id: ")
		buf.WriteString(*e.id + "\n")
	}
	if e.name != nil {
		buf.WriteString("event: ")
		buf.WriteString(*e.name + "\n")
	}
	if e.retry != nil {
		buf.WriteString("retry: ")
		buf.WriteString(*e.retry + "\n")
	}
	if e.data != nil {
		buf.WriteString("data: ")
		bs, err := json.Marshal(e.data)
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

// Send syntax sugar
func (e Event) Send(w io.Writer, f http.Flusher) error {
	_, err := w.Write(e.Bytes())
	if err != nil {
		return err
	}
	f.Flush()
	return nil
}

// NewEvent creates a new SSE event.
func NewEvent(name string, data any) Event {
	return Event{
		name: &name,
		data: data,
	}
}

// NewMessage is a special Event using a  single/default message in an SSE stream.
func NewMessage(data any) Event {
	return Event{
		data: data,
	}
}
