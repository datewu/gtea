package sse

import (
	"bytes"
	"encoding/json"
)

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
