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
