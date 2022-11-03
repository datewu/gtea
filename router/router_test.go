package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datewu/gtea/handler"
)

func reqTestHelper(method, path string, body io.Reader,
	h http.Handler, code int, expect string, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	h.ServeHTTP(w, req)
	if w.Code != code {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != expect {
		t.Errorf("expected %q got %q", expect, w.Body.String())
	}
}

func getReqHelper(path string, h http.Handler, code int, expect string, t *testing.T) {
	reqTestHelper(http.MethodGet, path, nil, h, code, expect, t)
}

func TestRouterHealhCheck(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	expect := `{"error":"the requested resource could not be found"}`
	getReqHelper("/v1/healthcheck", r.Handler(), http.StatusNotFound, expect, t)
}

func TestSuffix(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)

	r.Get("/health", handler.HealthCheck)
	r.Get("/abc/", handler.HealthCheck)
	r.Get("/", handler.HealthCheck)
	expect := `{"status":"available"}`
	h := r.Handler()
	getReqHelper("/health", h, http.StatusOK, expect, t)
	getReqHelper("/health/", h, http.StatusOK, expect, t)

	getReqHelper("/abc/", h, http.StatusOK, expect, t)
	getReqHelper("/abc", h, http.StatusOK, expect, t)

	getReqHelper("/", h, http.StatusOK, expect, t)
	getReqHelper("/a", h, http.StatusNotFound, `{"error":"the requested resource could not be found"}`, t)
}
func TestRequestMethods(t *testing.T) {
	conf := &Config{}
	ro := NewRouter(conf)

	ro.Get("/ok", handler.HealthCheck)
	ro.Post("/ok", handler.HealthCheck)
	ro.Put("/ok", handler.HealthCheck)
	ro.Delete("/ok", handler.HealthCheck)
	h := ro.Handler()
	request := func(method string) {
		expect := `{"status":"available"}`
		reqTestHelper(method, "/ok", nil, h,
			http.StatusOK, expect, t)
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
			"name": handler.ReadPathParam(r, "name"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	nameCityHandle := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"name": handler.ReadPathParam(r, "name"),
			"city": handler.ReadPathParam(r, "city"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	locationHandle := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"country": handler.ReadPathParam(r, "country"),
			"city":    handler.ReadPathParam(r, "city"),
		}
		handler.WriteJSON(w, http.StatusOK, data, nil)
	}
	ro.Get("/hi/:name", nameHandle)
	ro.Get("/hi/:name/:city", nameCityHandle)
	ro.Get("/hi/:country/:city/good", locationHandle)
	h := ro.Handler()
	nameReq := func(name string) {
		expect := fmt.Sprintf(`{"name":"%s"}`, name)
		getReqHelper("/hi/"+name, h, http.StatusOK, expect, t)
	}
	nameCityReq := func(n, c string) {
		expect := fmt.Sprintf(`{"city":"%s","name":"%s"}`, c, n)
		getReqHelper("/hi/"+n+"/"+c, h, http.StatusOK, expect, t)
	}
	locationReq := func(c, city string) {
		expect := fmt.Sprintf(`{"city":"%s","country":"%s"}`, city, c)
		getReqHelper("/hi/"+c+"/"+city+"/good", h, http.StatusOK, expect, t)
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

	expect := msg + `{"status":"available"}`
	getReqHelper("/", r.Handler(), http.StatusOK, expect, t)
}

func TestAggNormalMiddler(t *testing.T) {
	conf := &Config{}
	r := NewRouter(conf)
	_, ms, msg := newMiddlerwares(false)
	r.aggMiddleware = handler.AggregateMds(ms)
	r.Get("/", handler.HealthCheck)
	expect := msg + `{"status":"available"}`
	getReqHelper("/", r.Handler(), http.StatusOK, expect, t)
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
	expect := msg
	getReqHelper("/", r.Handler(), http.StatusOK, expect, t)
}
