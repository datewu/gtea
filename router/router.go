package router

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

// Router ..
type Router struct {
	conf          *Config
	trie          *pathTrie
	aggMiddleware handler.Middleware
}

// Handler serveHTTP
type Handler struct {
	trie pathTrie
	md   handler.Middleware
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var hf http.HandlerFunc
	tHandler := h.trie.get(r.Method + r.URL.Path)
	if tHandler == nil {
		hf = handler.NotFoundMsg("the requested resource could not be found")
	} else {
		hf = tHandler.ServeHTTP
	}
	h.md(hf)(w, r)
}

func (ro *Router) Handler() Handler {
	if ro.conf.Debug {
		ro.trie.printPaths()
	}
	ro.aggBuildInMiddlewares()
	h := Handler{
		trie: *ro.trie,
		md:   ro.aggMiddleware,
	}
	if h.md == nil {
		h.md = handler.VoidMiddleware
	}
	return h
}

func NewRouter(c *Config) *Router {
	r := &Router{conf: c}
	r.trie = newPathTrie()
	return r
}

func (r *Router) Handle(method, path string, h http.Handler) {
	r.trie.put(method+path, h)
}

func (r *Router) HandleFunc(method, path string, hf http.HandlerFunc) {
	if hf == nil {
		panic("http: nil http.handlerFunc")
	}
	r.Handle(method, path, hf)
}

func (r *Router) ServeFiles(path string, root http.Dir) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	r.trie.putEnd(http.MethodGet+path, h)
}

func (r *Router) ServeFilesWithGzip(path string, root http.Dir) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	hf := handler.GzipMiddleware(h.ServeHTTP)
	r.trie.putEnd(http.MethodGet+path, hf)
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

// Static serve dir dest
func (r *Router) Static(path string, dst string) {
	r.ServeFiles(path, http.Dir(dst))
}

// StaticGZIP serve dir dest with Gzip middleware
func (r *Router) StaticGZIP(path string, dst string) {
	r.ServeFilesWithGzip(path, http.Dir(dst))
}
