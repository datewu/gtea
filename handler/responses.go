package handler

import (
	"net/http"
	"time"
)

// SetSimpleCookie set key value within a week
func SetSimpleCookie(w http.ResponseWriter, r *http.Request, k, v string) {
	du := 7 * 24 * time.Hour
	expire := time.Now().Add(du)
	cookie := http.Cookie{
		Name: k, Value: v,
		Path:    "/",
		Domain:  r.URL.Host,
		Expires: expire, MaxAge: int(du.Seconds()),
	}
	http.SetCookie(w, &cookie)
}

// OKJSON response 200 respose with a json data
func OKJSON(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, data, nil)
}

// OKText response 200 with the plain text
func OKText(w http.ResponseWriter, text string) {
	WriteStr(w, http.StatusOK, text, nil)
}

// NotFound general 404 response
func NotFound(w http.ResponseWriter) {
	errResponse(w, http.StatusNotFound,
		"the requested resource could not be found")
}

// EditConflict 409 response
func EditConflict(w http.ResponseWriter) {
	errResponse(w, http.StatusConflict,
		"unable to update the record due to an edit conflict, please try later")
}

// RateLimitExceede  429 response
func RateLimitExceede(w http.ResponseWriter) {
	errResponse(w, http.StatusTooManyRequests,
		"rate limit exceeded")
}

// InvalidCredentials 400 response for bad reruest
func InvalidCredentials(w http.ResponseWriter) {
	errResponse(w, http.StatusBadRequest,
		"invalid authentication credentials")
}

// InvalidAuthenticationToken 401 response
func InvalidAuthenticationToken(w http.ResponseWriter) {
	errResponse(w, http.StatusUnauthorized,
		"invalid or missing authentication token")
}

// AuthenticationRequire 401 response tip
func AuthenticationRequire(w http.ResponseWriter) {
	errResponse(w, http.StatusUnauthorized,
		"you must be authenticated to access this resource")
}

// InactiveAccount 403 response for inactionaccount
func InactiveAccount(w http.ResponseWriter) {
	errResponse(w, http.StatusForbidden,
		"your user account must be activated to access this resource")
}

// NotPermitted 403 response
func NotPermitted(w http.ResponseWriter) {
	errResponse(w, http.StatusForbidden,
		"your user account doesn't have the necessary permissions to access this resource")
}

// MethodNotAllow 405 response
func MethodNotAllow(w http.ResponseWriter) {
	MethodNotAllowed(w, nil)
}

// HandleBadRequest  400 response with custom message
func BadRequestMsg(w http.ResponseWriter, msg string) {
	errResponse(w, http.StatusBadRequest, msg)
}

// BadRequestErr 400 response with a error
func BadRequestErr(w http.ResponseWriter, err error) {
	BadRequestMsg(w, err.Error())
}

// FailedValidation 400 response
func FailedValidation(w http.ResponseWriter, errs map[string]string) {
	errResponse(w, http.StatusBadRequest, errs)
}

// ServerErr a general 500 response with an err
func ServerErr(w http.ResponseWriter, err error) {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err.Error(),
	}
	errResponse(w, http.StatusInternalServerError, errs)
}

// ServerErrAny a general 500 response
func ServerErrAny(w http.ResponseWriter, msg interface{}) {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": msg,
	}
	errResponse(w, http.StatusInternalServerError, errs)
}
