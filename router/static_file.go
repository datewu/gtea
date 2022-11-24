package router

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

func (r *Router) ServeFS(path string, root http.FileSystem, mds ...handler.Middleware) {
	fs := http.FileServer(root)
	h := http.StripPrefix(path, fs)
	path = strings.TrimSuffix(path, "/")
	for _, m := range mds {
		if m == nil {
			continue
		}
		h = m(h.ServeHTTP)
	}
	r.trie.putEnd(http.MethodGet+path, h)
}

func (r *Router) ServeFSWithGzip(path string, root http.FileSystem, mds ...handler.Middleware) {
	mds = append(mds, handler.GzipMiddleware)
	r.ServeFS(path, root, mds...)
}

// Static serve dir dest
func (r *Router) Static(path string, dst string, mds ...handler.Middleware) {
	r.ServeFS(path, http.Dir(dst), mds...)
}

// StaticGZIP serve dir dest with Gzip middleware
func (r *Router) StaticGZIP(path string, dst string, mds ...handler.Middleware) {
	r.ServeFSWithGzip(path, http.Dir(dst), mds...)
}

// Static serve dir dest
func (g *RoutesGroup) Static(path string, dst string) {
	g.r.Static(g.prefix+path, dst, g.middlerware)
}

// StaticGZIP serve dir dest with Gzip middleware
func (g *RoutesGroup) StaticGZIP(path string, dst string) {
	g.r.ServeFSWithGzip(g.prefix+path, http.Dir(dst), g.middlerware)
}

// ServerFs ...
func (g *RoutesGroup) ServeFS(path string, root http.FileSystem) {
	g.r.ServeFS(g.prefix+path, root, g.middlerware)
}

// ServFSWithGzip ...
func (g *RoutesGroup) ServeFSWithGzip(path string, root http.FileSystem) {
	g.r.ServeFSWithGzip(g.prefix+path, root, g.middlerware)
}
