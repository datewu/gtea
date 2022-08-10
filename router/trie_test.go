package router

import (
	"net/http"
	"testing"
)

type dumpHandler int

// ServeHTTP impl http.Handler
func (d dumpHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestNormalTrie(t *testing.T) {
	p := newPathTrie()
	table := map[string]http.Handler{
		"get":                dumpHandler(1),
		"get/b":              dumpHandler(10),
		"get/a/bob/cow":      dumpHandler(2),
		"get/a/bob/cow/lol":  dumpHandler(3),
		"post":               dumpHandler(4),
		"post/c":             dumpHandler(40),
		"post/a/bob/cow":     dumpHandler(5),
		"post/a/bob/cow/lol": dumpHandler(6),
		"delete/x/y/z/ko":    dumpHandler(11),
		"delete/x/y/z":       dumpHandler(12),
		"delete/x/y":         dumpHandler(13),
		"delete/x":           dumpHandler(14),
		"delete/":            dumpHandler(15),
		"delete":             dumpHandler(16),
	}
	for k, v := range table {
		p.put(k, v)
		res := p.get(k)
		if res != v {
			t.Error("key:", k, "want:", v, "got", res)
		}
	}
	nilRes := p.get("noexist")
	if nilRes != nil {
		t.Error("should return nil, got:", nilRes)
	}
}

func TestRegTrie(t *testing.T) {
	p := newPathTrie()
	putTable := map[string]http.Handler{
		"get/:name/:age":       dumpHandler(1),
		"get/:name/:age/hello": dumpHandler(2),
		"get/a":                dumpHandler(3),
		"post/:c":              dumpHandler(4),
		"get/:one/:two/three":  dumpHandler(5),
	}
	getTable := map[string]http.Handler{
		"get/bob/11":             putTable["get/:name/:age"],
		"get/lol/22/hello":       putTable["get/:name/:age/hello"],
		"get/lol/22/hello/":      nil,
		"get/lol/22/hello_u":     nil,
		"get/a":                  putTable["get/a"],
		"get/abc":                nil,
		"post/:c":                putTable["post/:c"],
		"post/lol":               putTable["post/:c"],
		"get/foo/baro/three":     putTable["get/:one/:two/three"],
		"get/foo/baro/three/":    nil,
		"get/foo/baro/three_lol": nil,
	}
	for k, v := range putTable {
		p.put(k, v)
	}
	for k, v := range getTable {
		res := p.get(k)
		if res != v {
			t.Error("key:", k, "want:", v, "got", res)
		}
	}
}
