package handler

import (
	"encoding/json"
	"net/http"
)

func errResponse(w http.ResponseWriter, code int, msg any) {
	data := Envelope{"error": msg}
	WriteJSON(w, code, data, nil)
}

// Envelope is a JSON envelope for better client response
type Envelope map[string]any

// WriteStr writes a string to the response.
func WriteStr(w http.ResponseWriter, status int, msg string, headers http.Header) {
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

// writeJSON writes a JSON object to the response.
func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) {
	js, err := json.Marshal(data)
	if err != nil {
		msg := Envelope{"error": err}
		WriteJSON(w, http.StatusInternalServerError, msg, nil)
		return
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}
