package handler

import (
	"net/http"
)

// HandleHelper is a handler for all requests
type HandleHelper struct {
	w http.ResponseWriter
	r *http.Request
}

// NewHandleHelper returns a new handler
func NewHandleHelper(w http.ResponseWriter, r *http.Request) *HandleHelper {
	return &HandleHelper{
		w: w,
		r: r,
	}
}

// NotFound handle 404 response
func (h *HandleHelper) NotFound() {
	errResponse(http.StatusNotFound,
		"the requested resource could not be found",
	)(h.w, h.r)
}

// EditConflict handle 409 response
func (h *HandleHelper) EditConflict() {
	errResponse(http.StatusConflict,
		"unable to update the record due to an edit conflict, please try later",
	)(h.w, h.r)
}

// RateLimitExceede handle 429 response
func (h *HandleHelper) RateLimitExceede() {
	errResponse(http.StatusTooManyRequests,
		"rate limit exceeded",
	)(h.w, h.r)
}

// InvalidCredentials handle 400 response
func (h *HandleHelper) InvalidCredentials() {
	errResponse(http.StatusBadRequest,
		"invalid authentication credentials",
	)(h.w, h.r)
}

// InvalidAuthenticationToken handle 401 response
func (h *HandleHelper) InvalidAuthenticationToken() {
	errResponse(http.StatusUnauthorized,
		"invalid or missing authentication token",
	)(h.w, h.r)
}

// AuthenticationRequire handle 401 response
func (h *HandleHelper) AuthenticationRequire() {
	errResponse(http.StatusUnauthorized,
		"you must be authenticated to access this resource",
	)(h.w, h.r)
}

// InactiveAccount handle 403 response
func (h *HandleHelper) InactiveAccount() {
	errResponse(http.StatusForbidden,
		"your user account must be activated to access this resource",
	)(h.w, h.r)
}

// NotPermitted handle 403 response
func (h *HandleHelper) NotPermitted() {
	errResponse(http.StatusForbidden,
		"your user account doesn't have the necessary permissions to access this resource",
	)(h.w, h.r)
}

// MethodNotAllow handle 405 response
func (h *HandleHelper) MethodNotAllow() {
	MethodNotAllowed(h.w, h.r)
}

// HandleBadRequest handle 400 response with custom message
func (h *HandleHelper) BadRequestMsg(msg string) {
	errResponse(http.StatusBadRequest, msg)(h.w, h.r)
}

// BadRequestErr handle 400 response with a error
func (h *HandleHelper) BadRequestErr(err error) {
	h.BadRequestMsg(err.Error())
}

// FailedValidation handle 400 response
func (h *HandleHelper) FailedValidation(errs map[string]string) {
	errResponse(http.StatusBadRequest, errs)(h.w, h.r)
}

// ServerErr handle 500 response
func (h *HandleHelper) ServerErr(err error) {
	errs := map[string]interface{}{
		"error":  "the server encountered a problem and could not process your request",
		"detail": err.Error(),
	}
	errResponse(http.StatusInternalServerError, errs)(h.w, h.r)
}
