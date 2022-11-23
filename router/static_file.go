package router

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

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

// Static serve dir dest
func (r *Router) Static(path string, dst string) {
	r.ServeFiles(path, http.Dir(dst))
}

// StaticGZIP serve dir dest with Gzip middleware
func (r *Router) StaticGZIP(path string, dst string) {
	r.ServeFilesWithGzip(path, http.Dir(dst))
}
