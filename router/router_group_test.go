package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea/handler"
)

func TestUseMiddlerware(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	_, ms, msg := newMiddlerwares(false)
	g.Use(ms...)
	buildInHealthEscapeAggMiddles := func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	othersApplyMiddlewares := func() {
		g.Get("/abc", handler.HealthCheck)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/abc", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := msg + `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
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
	buildInHealthEscapeAggMiddles := func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	othersApplyMiddlewares := func() {
		g.Get("/abc", handler.HealthCheck)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/abc", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := msg
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
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
	buildInHealthEscapeAggMiddles := func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		hg.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	othersApplyMiddlewares := func() {
		hg.Get("/abc", handler.HealthCheck)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/hello/abc", nil)
		hg.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := msg + `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := `{"status":"available"}`
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}

func TestGroup(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	a := g.Group("/a")
	b := g.Group("/b")
	a.Get("/ok", handler.HealthCheck)
	b.Get("/ok", handler.HealthCheck)
	okPath := func(path string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	okPath("/a/ok")
	okPath("/b/ok")

	notOKPath := func(path string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected %d got %d", http.StatusNotFound, w.Code)
		}
		expect := `{"error":"the requested resource could not be found"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	notOKPath("/a/notok")
	notOKPath("/c/ok")
}
