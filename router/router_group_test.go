package router

import (
	"net/http"
	"testing"

	"github.com/datewu/gtea/handler"
)

func defaultHealthcheckHelper(h http.Handler, t *testing.T) {
	expect := `{"status":"available"}`
	getReqHelper("/v1/healthcheck", h, http.StatusOK, expect, t)
}

func TestUseMiddlerware(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	_, ms, msg := newMiddlerwares(false)
	g.Use(ms...)
	g.Get("/abc", handler.HealthCheck)
	h := g.Handler()
	buildInHealthEscapeAggMiddles := func() {
		defaultHealthcheckHelper(h, t)
	}
	othersApplyMiddlewares := func() {
		expect := msg + `{"status":"available"}`
		getReqHelper("/abc", h, http.StatusOK, expect, t)
	}
	buildInHealthEscapeAggMiddles()
	othersApplyMiddlewares()
}

func TestUseMiddlerwareWithAbort(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	msgs, ms, msg := newMiddlerwares(true)
	if len(ms) != len(msgs) {
		t.Errorf("expected %d middlers got %d", len(msgs), len(ms))
	}
	g.Use(ms...)
	g.Get("/abc", handler.HealthCheck)
	h := g.Handler()
	buildInHealthEscapeAggMiddles := func() {
		defaultHealthcheckHelper(h, t)
	}
	othersApplyMiddlewares := func() {
		getReqHelper("/abc", h, http.StatusOK, msg, t)
	}
	buildInHealthEscapeAggMiddles()
	othersApplyMiddlewares()
}

func TestGrupWithMiddlerware(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	_, ms, msg := newMiddlerwares(false)
	hg := g.Group("/hello", ms...)
	hg.Get("/abc", handler.HealthCheck)
	h := hg.Handler()
	buildInHealthEscapeAggMiddles := func() {
		defaultHealthcheckHelper(h, t)
	}
	othersApplyMiddlewares := func() {
		expect := msg + `{"status":"available"}`
		getReqHelper("/hello/abc", h, http.StatusOK, expect, t)
	}
	buildInHealthEscapeAggMiddles()
	othersApplyMiddlewares()
}

func TestGroupHealhCheck(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	defaultHealthcheckHelper(g.Handler(), t)
}

func TestGroupRequestMethods(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	g.Get("/ok", handler.HealthCheck)
	g.Post("/ok", handler.HealthCheck)
	g.Put("/ok", handler.HealthCheck)
	g.Delete("/ok", handler.HealthCheck)
	h := g.Handler()
	expect := `{"status":"available"}`
	request := func(method string) {
		reqTestHelper(method, "/ok", nil, h,
			http.StatusOK, expect, t)
	}
	request(http.MethodGet)
	request(http.MethodPost)
	request(http.MethodPut)
	request(http.MethodDelete)
}

func TestGroup(t *testing.T) {
	rconf := &Config{
		Metrics: true,
	}
	rconf.Limiter.Burst = 10
	rconf.Limiter.Rps = 100
	rconf.Limiter.Enabled = true
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	a := g.Group("/a")
	b := g.Group("/b")
	a.Get("/ok", handler.HealthCheck)
	b.Get("/ok", handler.HealthCheck)
	h := g.Handler()
	okPath := func(path string) {
		expect := `{"status":"available"}`
		getReqHelper(path, h, http.StatusOK, expect, t)
	}
	okPath("/a/ok")
	okPath("/b/ok")

	notOKPath := func(path string) {
		expect := `{"error":"the requested resource could not be found"}`
		getReqHelper(path, h, http.StatusNotFound, expect, t)
	}
	notOKPath("/a/notok")
	notOKPath("/c/ok")

	c := g.Group("/api/v1")
	c.Get("/", handler.HealthCheck)
	okPath("/api/v1")
	okPath("/api/v1/")
}
