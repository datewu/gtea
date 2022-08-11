package router

import (
	"errors"
	"net/http"
	"strconv"
)

// ReadPath returns the string param value in the request path
func ReadPath(r *http.Request, name string) string {
	keys := r.Context().Value(paramsCtxKey)
	ks, ok := keys.([]string)
	if !ok {
		return ""
	}
	values := r.Context().Value(paramsCtxValue)
	vs, ok := values.([]string)
	if !ok {
		return ""
	}
	if len(ks) != len(vs) {
		return ""
	}
	for i, v := range ks {
		if v == name {
			return vs[i]
		}
	}
	return ""
}

// ReadInt64Path returns the int64 param value in the request path
func ReadInt64Path(r *http.Request, name string) (int64, error) {
	v := ReadPath(r, name)
	if v == "" {
		return 0, errors.New("empty param")
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
