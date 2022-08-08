package router

import (
	"net/http"
)

// bag holds all paths relative funcs
type bag struct {
	rt     *Router
	config *Config
}

// Router ..
type Router struct {
	mux                        *http.ServeMux
	NotFound, MethodNotAllowed http.HandlerFunc
}

func (ro *Router) interceptHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		i := &interceptor4xx{
			origWriter:       w,
			methodNotAllowed: ro.MethodNotAllowed,
			notFound:         ro.NotFound,
		}
		r.URL.Path += r.Method
		ro.mux.ServeHTTP(i, r)
	}
	return fn
}

func NewRouter() *Router {
	r := &Router{}
	r.mux = http.NewServeMux()
	r.NotFound = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	r.MethodNotAllowed = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return r
}

func (r *Router) Handle(method, path string, h http.Handler) {
	if !containerPathParam(path) {
		r.mux.Handle(path+method, h)
		return
	}
}

func (r *Router) HandleFunc(method, path string, hf http.HandlerFunc) {
	if hf == nil {
		panic("http: nil handler")
	}
	r.Handle(method, path, hf)
}

func (ro *Router) ServeFiles(path string, root http.Dir) {

}
