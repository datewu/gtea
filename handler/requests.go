package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ErrNoToken is returned when a token is not found in the request
var ErrNoToken = errors.New("no token")

func parseFormfile(r *http.Request, name string) (string, io.ReadCloser, error) {
	const maxMemory = 32 << 20 // 32 MB
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			return "", nil, err
		}
	}
	f, fh, err := r.FormFile(name)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	src, err := fh.Open()
	if err != nil {
		return "", nil, err
	}
	return fh.Filename, src, nil

}

// SaveFormFile write file to a dir with upload filename.
func SaveFormFile(r *http.Request, name, dir string) (string, error) {
	fn, src, err := parseFormfile(r, name)
	if err != nil {
		return "", err
	}
	defer src.Close()
	fullPath := filepath.Join(dir, fn)
	dst, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

// WriteFormFile write file to a io.Writer.
func WriteFormFile(r *http.Request, name string, dst io.Writer) error {
	_, src, err := parseFormfile(r, name)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	return err
}

// GetBearerToken returns the bearer token from the request
func GetBeareToken(r *http.Request, name string) (string, error) {
	head, err := GetToken(r, name)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(head, "Bearer ") {
		return "", errors.New("token must be a Bearer token")
	}
	return strings.TrimPrefix(head, "Bearer "), nil
}

// GetToken returns the token from the request
func GetToken(r *http.Request, name string) (string, error) {
	if name == "" {
		name = "token"
	}
	q := ReadQuery(r, name, "") // for ws query
	if q != "" {
		return q, nil
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", ErrNoToken
	}
	return token, nil
}

// GetValue gets a value from the request context.
func GetValue(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// SetValue sets a value on the returned request.
func SetValue(r *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

// PathRegs for ctx key
type PathRegs string

const (
	ParamsCtxKey   PathRegs = "path_param_names"
	ParamsCtxValue PathRegs = "path_param_values"
)

// ReadPathParam returns the string param value in the request path
func ReadPathParam(r *http.Request, name string) string {
	keys := r.Context().Value(ParamsCtxKey)
	ks, ok := keys.([]string)
	if !ok {
		return ""
	}
	values := r.Context().Value(ParamsCtxValue)
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

// ReadInt64PathParam returns the int64 param value in the request path
func ReadInt64PathParam(r *http.Request, name string) (int64, error) {
	v := ReadPathParam(r, name)
	if v == "" {
		return 0, errors.New("empty param")
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// ReadQuery returns the string query value with a defaut value from the request
func ReadQuery(r *http.Request, key string, defaultValue string) string {
	qs := r.URL.Query()
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// ReadCSVQuery returns the csv list query value with a defaut list from the request
func ReadCSVQuery(r *http.Request, key string, defaultValue []string) []string {
	qs := r.URL.Query()
	cs := qs.Get(key)
	if cs == "" {
		return defaultValue
	}
	return strings.Split(cs, ",")
}

// ReadInt64Query returns the int64 query value with a defaut value from the request
func ReadInt64Query(r *http.Request, key string, defaultValue int64) int64 {
	qs := r.URL.Query()
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// ReadJSON reads the request body up to the max size and unmarshal it to the given struct
func ReadMaxJSON(w http.ResponseWriter, r *http.Request, dst interface{}, max int64) error {
	if max == 0 {
		max = 8 * 1_048_576 // 8MB for max readJSON body
	}
	r.Body = http.MaxBytesReader(w, r.Body, max)
	err := decodeJSON(r.Body, dst)
	if err != nil {
		if err.Error() == "http: request body too large" {
			return fmt.Errorf("body must not be larger than %d bytes", max)
		}
	}
	return nil
}

// ReadJSON reads the request body and unmarshal it to the given struct
func ReadJSON(r *http.Request, dst interface{}) error {
	return decodeJSON(r.Body, dst)
}

func decodeJSON(r io.Reader, dst interface{}) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalErr):
			if unmarshalErr.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalErr.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalErr.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

			// an open issue at https://github.com/golang/go/issues/29035
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalErr):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single json value")
	}
	return nil
}
