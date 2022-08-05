package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ReadParam returns the string param value in the request path
func ReadParams(r *http.Request, name string) string {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName(name)
}

// ReadInt64Param returns the int64 param value in the request path
func ReadInt64Param(r *http.Request, name string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(name), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
