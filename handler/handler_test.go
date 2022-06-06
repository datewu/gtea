package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSetValue(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
	name := "old_man"
	age := 98

	r := SetValue(req, name, age)
	v := GetValue(r, name)
	a, ok := v.(int)
	if !ok {
		t.Fatalf("expected int type got %T", v)
	}
	if a != age {
		t.Fatalf("expected %d got %d", age, a)
	}
}

func TestHTTP(t *testing.T) {
	text := "hello, you"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	}))
	defer ts.Close()
	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != text {
		t.Fatalf("expected %q got %q", text, string(b))
	}
}

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
	a.Get("/ok", HealthCheck)
	b.Get("/ok", HealthCheck)
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
