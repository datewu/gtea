package router

import (
	"net/http"
	"strings"
	"sync"

	"github.com/datewu/gtea/handler"
)

// RoutesGroup is a group of routes
type RoutesGroup struct {
	r           *bag
	prefix      string
	middlewares []handler.Middleware
	once        sync.Once
	serverHTTP  http.HandlerFunc
}

func (g *RoutesGroup) bound() {
	g.r.rt.NotFound = handler.NotFoundMsg(
		"the requested resource could not be found")
	g.r.rt.MethodNotAllowed = handler.MethodNotAllowed
	g.Get("/v1/healthcheck", handler.HealthCheck)
	mm := g.r.rt.ServeHTTP
	middlewares := append(g.middlewares, g.r.buildIns()...)
	for _, m := range middlewares {
		mm = m(mm)
	}
	g.serverHTTP = mm
}

func (g *RoutesGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.once.Do(g.bound)
	g.serverHTTP(w, r)
}

// Use add middleware to the group
// middleware will be called in the order of use
// Or call NewHandler to add a new middleware
func (g *RoutesGroup) Use(mds ...handler.Middleware) {
	g.middlewares = append(g.middlewares, mds...)
}

// Group add a prefix to all path, for each Gropu call
// prefix will accumulate while middleware don't
func (g *RoutesGroup) Group(path string, mds ...handler.Middleware) *RoutesGroup {
	gp := &RoutesGroup{
		r:      g.r,
		prefix: g.prefix + path,
	}
	if mds != nil {
		gp.middlewares = append(gp.middlewares, mds...)
	}
	return gp
}

// HandleFunc handle new http request
func (g *RoutesGroup) HandleFunc(method, path string, handler http.HandlerFunc) {
	for _, v := range g.middlewares {
		handler = v(handler)
	}
	g.r.rt.HandleFunc(method, g.prefix+path, handler)
}

// Get is a shortcut for NewHandler(http.MethodGet, path, handler)
func (g *RoutesGroup) Get(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodGet, path, handler)
}

// Post is a shortcut for NewHandler(http.MethodPost, path, handler)
func (g *RoutesGroup) Post(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPost, path, handler)
}

// Put is a shortcut for NewHandler(http.MethodPut, path, handler)
func (g *RoutesGroup) Put(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPut, path, handler)
}

// Patch is a shortcut for NewHandler(http.MethodPatch, path, handler)
func (g *RoutesGroup) Patch(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodPatch, path, handler)
}

// Delete is a shortcut for NewHandler(http.MethodDelete, path, handler)
func (g *RoutesGroup) Delete(path string, handler http.HandlerFunc) {
	g.HandleFunc(http.MethodDelete, path, handler)
}

// Static is a shortcut for NewHandler(http.MethodDelete, path, handler)
func (g *RoutesGroup) Static(prefix string, dst string) {
	path := strings.TrimSuffix(prefix, "/") + "/*filepath"
	root := http.Dir(dst)
	g.r.rt.ServeFiles(path, root)
}
