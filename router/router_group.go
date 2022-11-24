package router

import (
	"errors"
	"net/http"

	"github.com/datewu/gtea/handler"
)

// RoutesGroup is a group of routes
type RoutesGroup struct {
	r           *Router
	prefix      string
	middlerware handler.Middleware
}

// NewRoutesGroup return a new routesgroup
func NewRoutesGroup(conf *Config) (*RoutesGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := NewRouter(conf)
	g := &RoutesGroup{r: r}
	g.Get("/v1/healthcheck", handler.HealthCheck)
	return g, nil
}

// Handler return http.Handler
func (g *RoutesGroup) Handler() Handler {
	return g.r.Handler()
}

// Use add middleware to the group
// middleware will be called in the order of use
// Or call NewHandler to add a new middleware
func (g *RoutesGroup) Use(mds ...handler.Middleware) {
	for _, v := range mds {
		g.middlerware = handler.Append(g.middlerware, v)
	}
}

// Group add a prefix to all path, for each Gropu call
// prefix and middleware will accumulate
func (g *RoutesGroup) Group(path string, mds ...handler.Middleware) *RoutesGroup {
	gp := &RoutesGroup{
		r:           g.r,
		middlerware: g.middlerware,
		prefix:      g.prefix + path,
	}
	for _, v := range mds {
		gp.middlerware = handler.Append(gp.middlerware, v)
	}
	return gp
}

// HandleFunc handle new http request
func (g *RoutesGroup) HandleFunc(method, path string, handler http.HandlerFunc) {
	if g.middlerware != nil {
		handler = g.middlerware(handler)
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
