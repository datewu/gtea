package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea/handler"
)

func TestRouterHealhCheck(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d got %d", http.StatusNotFound, w.Code)
	}
	expect := ""
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}

func TestRequestMethods(t *testing.T) {
	conf := &Config{}
	ro := NewRouter(conf)

	ro.Get("/ok", handler.HealthCheck)
	ro.Post("/ok", handler.HealthCheck)
	ro.Put("/ok", handler.HealthCheck)
	ro.Delete("/ok", handler.HealthCheck)
	request := func(method string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/ok", nil)
		ro.ServeHTTP(w, req)
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
	conf := &Config{}
	ro := NewRouter(conf)
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
	ro.Get("/hi/:name", nameHandle)
	ro.Get("/hi/:name/:city", nameCityHandle)
	ro.Get("/hi/:country/:city/good", locationHandle)
	nameReq := func(name string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hi/"+name, nil)
		ro.ServeHTTP(w, req)
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
		ro.ServeHTTP(w, req)
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
		ro.ServeHTTP(w, req)
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

func newNormalMiddler(injectMsg string) handler.Middleware {
	mid := func(next http.HandlerFunc) http.HandlerFunc {
		res := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(injectMsg))
			next(w, r)
		}
		return res
	}
	return mid
}

func newMiddlerwares(abort bool) ([]string, []handler.Middleware, string) {
	abortMsg := "abort followed http.handlers"
	msgs := []string{
		"inject msg 1 ",
		"inject msg 2 ",
		abortMsg,
		"inject msg 3 ",
		"inject msg 4 ",
		"inject msg 5 ",
	}
	if !abort {
		msgs[2] = "not abort"
	}
	ms := []handler.Middleware{}
	for _, v := range msgs {
		if v == abortMsg {
			ms = append(ms, handler.AbortMiddleware)
		} else {
			ms = append(ms, newNormalMiddler(v))
		}
	}
	msg := ""
	for _, v := range msgs {
		if v == abortMsg {
			break
		}
		msg += v
	}
	return msgs, ms, msg
}

func TestNormalMiddler(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	msg := " inject a msg "
	r.aggMiddleware = newNormalMiddler(msg)
	r.Get("/", handler.HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := msg + `{"status":"available"}`
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}

func TestAggNormalMiddler(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	_, ms, msg := newMiddlerwares(false)
	r.aggMiddleware = handler.AggregateMds(ms)
	r.Get("/", handler.HealthCheck)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := msg + `{"status":"available"}`
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}

func TestAggNormalMiddlerWithAbort(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	msgs, ms, msg := newMiddlerwares(true)
	if len(ms) != len(msgs) {
		t.Errorf("expected %d middlers got %d", len(msgs), len(ms))
	}
	r.aggMiddleware = handler.AggregateMds(ms)
	r.Get("/", handler.HealthCheck)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	expect := msg
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}
