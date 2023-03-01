package handler

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler a simple health check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"status": "available",
	}
	WriteJSON(w, http.StatusOK, data, nil)
}

// MethodNotAllowed method not found handler
var MethodNotAllowed http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	errResponse(w, http.StatusMethodNotAllowed, msg)
}

// NotFoundMsg method not found with custom message
func NotFoundMsg(msg string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		errResponse(w, http.StatusNotFound, msg)
	}
	return fn
}

// Ok is a high order func
func OkHandler(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, data, nil)
	}
}

// Version response version
func Version(ver string) http.HandlerFunc {
	data := Envelope{
		"version": ver,
	}
	return OkHandler(data)
}
