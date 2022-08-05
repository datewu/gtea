package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea/handler"
)

func TestHealhCheck(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
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

func TestPathParams(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	nameHandle := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"name": ReadParams(r, "name"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	g.Get("/hi/:name", nameHandle)
	namePath := func(name string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hi/"+name, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := fmt.Sprintf(`{"name":"%s"}`, name)
		if w.Body.String() != expect {
			t.Fatalf("expected %q got %q", expect, w.Body.String())
		}
	}
	namePath("joe-boy")
}
