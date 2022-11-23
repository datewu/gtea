package router

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

func (r *Router) ServeFS(path string, root http.FileSystem) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	r.trie.putEnd(http.MethodGet+path, h)
}

func (r *Router) ServeFSWithGzip(path string, root http.FileSystem) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	hf := handler.GzipMiddleware(h.ServeHTTP)
	r.trie.putEnd(http.MethodGet+path, hf)
}

// Static serve dir dest
func (r *Router) Static(path string, dst string) {
	r.ServeFS(path, http.Dir(dst))
}

// StaticGZIP serve dir dest with Gzip middleware
func (r *Router) StaticGZIP(path string, dst string) {
	r.ServeFSWithGzip(path, http.Dir(dst))
}

// Static serve dir dest
func (g *RoutesGroup) Static(path string, dst string) {
	g.r.Static(g.prefix+path, dst)
}

// StaticGZIP serve dir dest with Gzip middleware
func (g *RoutesGroup) StaticGZIP(path string, dst string) {
	g.r.ServeFSWithGzip(g.prefix+path, http.Dir(dst))
}

// ServerFs ...
func (g *RoutesGroup) ServeFS(path string, root http.FileSystem) {
	g.r.ServeFS(g.prefix+path, root)
}

// ServFSWithGzip ...
func (g *RoutesGroup) ServeFSWithGzip(path string, root http.FileSystem) {
	g.r.ServeFSWithGzip(g.prefix+path, root)
}
