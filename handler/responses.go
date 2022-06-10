package handler

import (
	"fmt"
	"net/http"
)

// OKJSON handle 200 respose
func OKJSON(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, data, nil)
}

// OKText handle 200 respose
func OKText(w http.ResponseWriter, text string) {
	WriteStr(w, http.StatusOK, text, nil)
}

func errResponse(code int, msg interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := Envelope{"error": msg}
		WriteJSON(w, code, data, nil)
	}
}

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
	errResponse(http.StatusMethodNotAllowed, msg)
}

func NotFoundMsg(msg string) http.HandlerFunc {
	return errResponse(http.StatusNotFound, msg)
}
