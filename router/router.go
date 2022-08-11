package router

import (
	"net/http"

	"github.com/datewu/gtea/handler"
)

// Router ..
type Router struct {
	conf                       *Config
	trie                       *pathTrie
	aggMiddleware              handler.Middleware
	NotFound, MethodNotAllowed http.HandlerFunc
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	innerest := func(iw http.ResponseWriter, ir *http.Request) {
		tHandler := ro.trie.get(ir.Method + ir.URL.Path)
		if tHandler != nil {
			tHandler.ServeHTTP(iw, ir)
			return
		}
		ro.NotFound(iw, ir)
	}
	if ro.aggMiddleware != nil {
		ro.aggMiddleware(innerest)(w, r)
		return
	}
	innerest(w, r)
}

func NewRouter(c *Config) *Router {
	r := &Router{conf: c}
	r.trie = newPathTrie()
	r.NotFound = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	r.aggBuildInMiddlewares()
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
