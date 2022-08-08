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
	req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := `{"status":"available"}`
	if w.Body.String() != expect {
		t.Fatalf("expected %q got %q", expect, w.Body.String())
	}
}

func TestRequestMethods(t *testing.T) {
	rconf := &Config{}
	g, err := NewRoutesGroup(rconf)
	if err != nil {
		t.Fatal(err)
	}
	g.Get("/ok", handler.HealthCheck)
	g.Post("/ok", handler.HealthCheck)
	g.Put("/ok", handler.HealthCheck)
	g.Delete("/ok", handler.HealthCheck)
	request := func(method string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/ok", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Fatalf("expected %q got %q", expect, w.Body.String())
		}
	}
	request(http.MethodGet)
	request(http.MethodPost)
	request(http.MethodPut)
	request(http.MethodDelete)
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
			"name": ReadPath(r, "name"),
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
