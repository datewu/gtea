package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Event a single event in SSE stream using json marshal encoding
type Event struct {
	id, name *string
	data     any
	retry    *string
}

func strPoint(name string, data *string) string {
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s\n", name, *data)
}

// RawBytes set marshal the event as a byte slice. Escape '\n' to '\r\n'
func (e Event) RawBytes() []byte {
	if e.data == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	buf.WriteString(strPoint("id", e.id))
	buf.WriteString(strPoint("event", e.name))
	buf.WriteString(strPoint("retry", e.retry))
	buf.WriteString("data: ")
	js := new(bytes.Buffer)
	renc := json.NewEncoder(js)
	renc.SetEscapeHTML(false)
	err := renc.Encode(e.data)
	if err != nil {
		return nil
	}
	bs := js.Bytes()
	for _, b := range bs {
		if b == '\n' {
			buf.WriteString("\r\n")
		} else {
			buf.WriteByte(b)
		}
	}
	buf.WriteString("\n\n")
	return buf.Bytes()
}

// Bytes marshal the event as a byte slice. Escape '\n' to '\r\n'
func (e Event) Bytes() []byte {
	if e.data == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	buf.WriteString(strPoint("id", e.id))
	buf.WriteString(strPoint("event", e.name))
	buf.WriteString(strPoint("retry", e.retry))
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

// SendRaw syntax sugar
func (e Event) SendRaw(w io.Writer, f http.Flusher) error {
	_, err := w.Write(e.RawBytes())
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
