package handler

import (
	"net/http"
	"strings"
	"sync"
)

// RouterGroup is a group of routes
type RouterGroup struct {
	r           *router
	prefix      string
	middlewares []Middleware
	once        sync.Once
	serverHTTP  http.HandlerFunc
}

func (g *RouterGroup) bound() {
	g.r.router.NotFound = errResponse(http.StatusNotFound,
		"the requested resource could not be found",
	)
	g.r.router.MethodNotAllowed = MethodNotAllowed
	g.Get("/v1/healthcheck", HealthCheck)
	h := http.Handler(g.r.router)
	mm := h.ServeHTTP
	middlewares := append(g.middlewares, g.r.buildIns()...)
	for _, m := range middlewares {
		mm = m(mm)
	}
	g.serverHTTP = mm
}

func (g *RouterGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.once.Do(g.bound)
	g.serverHTTP(w, r)
}

// Use add middleware to the group
// middleware will be called in the order of use
// Or call NewHandler to add a new middleware
func (g *RouterGroup) Use(mds ...Middleware) {
	g.middlewares = append(g.middlewares, mds...)
}

// Group add a prefix to all path, for each Gropu call
// prefix will accumulate while middleware don't
func (g *RouterGroup) Group(path string, mds ...Middleware) *RouterGroup {
	gp := &RouterGroup{
		r:      g.r,
		prefix: g.prefix + path,
	}
	if mds != nil {
		gp.middlewares = append(gp.middlewares, mds...)
	}
	return gp
}

// NewHandler handle new http request
func (g *RouterGroup) NewHandler(method, path string, handler http.HandlerFunc) {
	for _, v := range g.middlewares {
		handler = v(handler)
	}
	g.r.router.HandlerFunc(method, g.prefix+path, handler)
}

// Get is a shortcut for NewHandler(http.MethodGet, path, handler)
func (g *RouterGroup) Get(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodGet, path, handler)
}

// Post is a shortcut for NewHandler(http.MethodPost, path, handler)
func (g *RouterGroup) Post(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPost, path, handler)
}

// Put is a shortcut for NewHandler(http.MethodPut, path, handler)
func (g *RouterGroup) Put(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPut, path, handler)
}

// Patch is a shortcut for NewHandler(http.MethodPatch, path, handler)
func (g *RouterGroup) Patch(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPatch, path, handler)
}

// Delete is a shortcut for NewHandler(http.MethodDelete, path, handler)
func (g *RouterGroup) Delete(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodDelete, path, handler)
}

// Static is a shortcut for NewHandler(http.MethodDelete, path, handler)
func (g *RouterGroup) Static(prefix string, dst string) {
	path := strings.TrimSuffix(prefix, "/") + "/*filepath"
	root := http.Dir(dst)
	g.r.router.ServeFiles(path, root)
}
