package handler

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler is a handler for health check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"status": "available",
	}
	WriteJSON(w, http.StatusOK, data, nil)
}

// MethodNotAllowed is a handler for method not found
var MethodNotAllowed http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	errResponse(w, http.StatusMethodNotAllowed, msg)
}

func NotFoundMsg(msg string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		errResponse(w, http.StatusNotFound, msg)
	}
	return fn
}
