package router

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

// Router ..
type Router struct {
	conf                       *Config
	trie                       *pathTrie
	aggMiddleware              handler.Middleware
	NotFound, MethodNotAllowed http.HandlerFunc
}

// ServeHTTP ...
func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	last := func(iw http.ResponseWriter, ir *http.Request) {
		tHandler := ro.trie.get(ir.Method + ir.URL.Path)
		if tHandler != nil {
			tHandler.ServeHTTP(iw, ir)
			return
		}
		ro.NotFound(iw, ir)
	}
	if ro.aggMiddleware != nil {
		ro.aggMiddleware(last)(w, r)
		return
	}
	last(w, r)
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
	if r.conf.Debug {
		r.trie.walk("/", 1)
	}
	r.trie.put(method+path, h)
}

func (r *Router) HandleFunc(method, path string, hf http.HandlerFunc) {
	if hf == nil {
		panic("http: nil handler")
	}
	r.Handle(method, path, hf)
}

func (r *Router) ServeFiles(path string, root http.Dir) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	r.trie.putEnd(http.MethodGet+path, h)
}

// Get is a shortcut for HandleFunc(http.MethodGet, path, handler)
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodGet, path, handler)
}

// Post is a shortcut for HandleFunc(http.MethodPost, path, handler)
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodPost, path, handler)
}

// Put is a shortcut for HandleFunc(http.MethodPut, path, handler)
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodPut, path, handler)
}

// Patch is a shortcut for HandleFunc(http.MethodPatch, path, handler)
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodPatch, path, handler)
}

// Delete is a shortcut for HandleFunc(http.MethodDelete, path, handler)
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodDelete, path, handler)
}

// Static is a shortcut for HandleFunc(http.MethodDelete, path, handler)
func (r *Router) Static(path string, dst string) {
	r.ServeFiles(path, http.Dir(dst))
}
