package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea/handler"
)

func TestHealhCheck(t *testing.T) {
	rconf := &Config{}
	g, err := NewRouterGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := `{"status":"available"}`
	if w.Body.String() != expect {
		t.Fatalf("expected %q got %q", expect, w.Body.String())
	}
}

func TestGroup(t *testing.T) {
	rconf := &Config{}
	g, err := NewRouterGroup(rconf)
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
			t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Fatalf("expected %q got %q", expect, w.Body.String())
		}
	}
	okPath("/a/ok")
	okPath("/b/ok")

	notOKPath := func(path string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
		}
		expect := `{"error":"the requested resource could not be found"}`
		if w.Body.String() != expect {
			t.Fatalf("expected %q got %q", expect, w.Body.String())
		}
	}
	notOKPath("/a/notok")
	notOKPath("/c/ok")
}
