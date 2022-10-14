package router

import (
	"errors"
	"net/http"

	"github.com/datewu/gtea/handler"
)

// RoutesGroup is a group of routes
type RoutesGroup struct {
	r             *Router
	prefix        string
	aggMiddleware handler.Middleware
}

// NewRoutesGroup return a new routesgroup
func NewRoutesGroup(conf *Config) (*RoutesGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := NewRouter(conf)
	r.NotFound = handler.NotFoundMsg(
		"the requested resource could not be found")

	// this "/v1/healthcheck" escape all routerGroup.aggMiddleware
	r.HandleFunc(http.MethodGet, "/v1/healthcheck", handler.HealthCheck)
	return &RoutesGroup{r: r}, nil
}

func (g *RoutesGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.r.ServeHTTP(w, r)
}

// Use add middleware to the group
// middleware will be called in the order of use
// Or call NewHandler to add a new middleware
func (g *RoutesGroup) Use(mds ...handler.Middleware) {
	if len(mds) == 0 {
		return
	}
	ms := make([]handler.Middleware, len(mds)+1)
	ms[0] = g.aggMiddleware
	for i, v := range mds {
		ms[i+1] = v
	}
	g.aggMiddleware = handler.AggregateMds(ms)
}

// Group add a prefix to all path, for each Gropu call
// prefix will accumulate while middleware don't
func (g *RoutesGroup) Group(path string, mds ...handler.Middleware) *RoutesGroup {
	gp := &RoutesGroup{
		r:      g.r,
		prefix: g.prefix + path,
	}
	if mds != nil {
		gp.aggMiddleware = handler.AggregateMds(mds)
	}
	return gp
}

// HandleFunc handle new http request
func (g *RoutesGroup) HandleFunc(method, path string, handler http.HandlerFunc) {
	if g.aggMiddleware != nil {
		handler = g.aggMiddleware(handler)
	}
	g.r.HandleFunc(method, g.prefix+path, handler)
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

// Static serve dir dest
func (g *RoutesGroup) Static(path string, dst string) {
	g.r.Static(path, dst)
}

// StaticGZIP serve dir dest with Gzip middleware
func (g *RoutesGroup) StaticGZIP(path string, dst string) {
	g.r.ServeFilesWithGzip(path, http.Dir(dst))
}
