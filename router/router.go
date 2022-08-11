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
	trie                       *pathTrie
	NotFound, MethodNotAllowed http.HandlerFunc
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tHandler := ro.trie.get(r.Method + r.URL.Path)
	if tHandler != nil {
		tHandler.ServeHTTP(w, r)
		return
	}
	ro.NotFound(w, r)
}

func NewRouter() *Router {
	r := &Router{}
	r.trie = newPathTrie()
	r.NotFound = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	return r
}

func (r *Router) Handle(method, path string, h http.Handler) {
	r.trie.put(method+path, h)
}

func (r *Router) HandleFunc(method, path string, hf http.HandlerFunc) {
	if hf == nil {
		panic("http: nil handler")
	}
	r.Handle(method, path, hf)
}

func (ro *Router) ServeFiles(path string, root http.Dir) {
	// TODO

}
