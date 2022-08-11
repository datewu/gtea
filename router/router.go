package router

import (
	"net/http"

	"github.com/datewu/gtea/handler"
)

// Router ..
type Router struct {
	conf                       *Config
	trie                       *pathTrie
	middleware                 handler.Middleware
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
	if ro.middleware != nil {
		ro.middleware(innerest)(w, r)
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
	ms := r.buildIns()
	if len(ms) > 0 {
		m := handler.VoidHandlerFunc
		for i := len(ms) - 1; i >= 0; i-- {
			m = ms[i](m)
		}
		// r.middleware =
	}

	// 	mm := g.r.ServeHTTP
	// 	middlewares := append(g.middlewares, g.r.buildIns()...)
	// 	for _, m := range middlewares {
	// 		mm = m(mm)
	// 	}
	// 	g.serverHTTP = mm
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
