package router

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

// ReadPath returns the string param value in the request path
func ReadPath(r *http.Request, name string) string {
	return "todo"
}

// ReadInt64Path returns the int64 param value in the request path
func ReadInt64Path(r *http.Request, name string) (int64, error) {
	v := ReadPath(r, name)
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

var pathReg = regexp.MustCompile(`/:(\w+)`)

func containerPathParam(path string) bool {
	return pathReg.MatchString(path)
}

func findPathParam(path string) []string {
	res := pathReg.FindAllStringSubmatch(path, -1)
	if res == nil {
		return nil
	}
	names := make([]string, len(res))
	for i := range res {
		names[i] = res[i][1]
	}
	return names
}
