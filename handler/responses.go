package handler

import (
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

// NotFound handle 404 response
func NotFound(w http.ResponseWriter) {
	errResponse(w, http.StatusNotFound,
		"the requested resource could not be found")
}

// EditConflict handle 409 response
func EditConflict(w http.ResponseWriter) {
	errResponse(w, http.StatusConflict,
		"unable to update the record due to an edit conflict, please try later")
}

// RateLimitExceede handle 429 response
func RateLimitExceede(w http.ResponseWriter) {
	errResponse(w, http.StatusTooManyRequests,
		"rate limit exceeded")
}

// InvalidCredentials handle 400 response
func InvalidCredentials(w http.ResponseWriter) {
	errResponse(w, http.StatusBadRequest,
		"invalid authentication credentials")
}

// InvalidAuthenticationToken handle 401 response
func InvalidAuthenticationToken(w http.ResponseWriter) {
	errResponse(w, http.StatusUnauthorized,
		"invalid or missing authentication token")
}

// AuthenticationRequire handle 401 response
func AuthenticationRequire(w http.ResponseWriter) {
	errResponse(w, http.StatusUnauthorized,
		"you must be authenticated to access this resource")
}

// InactiveAccount handle 403 response
func InactiveAccount(w http.ResponseWriter) {
	errResponse(w, http.StatusForbidden,
		"your user account must be activated to access this resource")
}

// NotPermitted handle 403 response
func NotPermitted(w http.ResponseWriter) {
	errResponse(w, http.StatusForbidden,
		"your user account doesn't have the necessary permissions to access this resource")
}

// MethodNotAllow handle 405 response
func MethodNotAllow(w http.ResponseWriter) {
	MethodNotAllowed(w, nil)
}

// HandleBadRequest handle 400 response with custom message
func BadRequestMsg(w http.ResponseWriter, msg string) {
	errResponse(w, http.StatusBadRequest, msg)
}

// BadRequestErr handle 400 response with a error
func BadRequestErr(w http.ResponseWriter, err error) {
	BadRequestMsg(w, err.Error())
}

// FailedValidation handle 400 response
func FailedValidation(w http.ResponseWriter, errs map[string]string) {
	errResponse(w, http.StatusBadRequest, errs)
}

// ServerErr handle 500 response
func ServerErr(w http.ResponseWriter, err error) {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err.Error(),
	}
	errResponse(w, http.StatusInternalServerError, errs)
}
