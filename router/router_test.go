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
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := `{"status":"available"}`
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
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
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := `{"status":"available"}`
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	request(http.MethodGet)
	request(http.MethodPost)
	request(http.MethodPut)
	request(http.MethodDelete)
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
	nameCityHandle := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"name": ReadPath(r, "name"),
			"city": ReadPath(r, "city"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	locationHandle := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"country": ReadPath(r, "country"),
			"city":    ReadPath(r, "city"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	g.Get("/hi/:name", nameHandle)
	g.Get("/hi/:name/:city", nameCityHandle)
	g.Get("/hi/:country/:city/good", locationHandle)
	nameReq := func(name string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hi/"+name, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := fmt.Sprintf(`{"name":"%s"}`, name)
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	nameCityReq := func(n, c string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hi/"+n+"/"+c, nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := fmt.Sprintf(`{"city":"%s","name":"%s"}`, c, n)
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	locationReq := func(c, city string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hi/"+c+"/"+city+"/good", nil)
		g.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		expect := fmt.Sprintf(`{"city":"%s","country":"%s"}`, city, c)
		if w.Body.String() != expect {
			t.Errorf("expected %q got %q", expect, w.Body.String())
		}
	}
	nameReq("joe-boy")
	nameCityReq("luffe", "dao-lol")
	locationReq("china", "hubei-wuhan")
}
