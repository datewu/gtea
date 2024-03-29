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

func TestTriewalk(t *testing.T) {
	p := newPathTrie()
	table := map[string]http.Handler{
		"/g/et":                   dumpHandler(1),
		"/get/b":                  dumpHandler(10),
		"/get/a/bob/cow":          dumpHandler(2),
		"/get/a/bob/:cow/lol":     dumpHandler(3),
		"/post":                   dumpHandler(4),
		"/post/:c":                dumpHandler(40),
		"/post/a/bob/cow":         dumpHandler(5),
		"/post/a/bob/cow/lol":     dumpHandler(6),
		"/delete/:x/y/z/ko":       dumpHandler(11),
		"/delete/x/y/z":           dumpHandler(12),
		"/delete/x/y":             dumpHandler(13),
		"/delete/x":               dumpHandler(14),
		"/delete":                 dumpHandler(16),
		"/hi/:name":               dumpHandler(33),
		"/hi/:name/:city":         dumpHandler(34),
		"/hi/:country/:city/good": dumpHandler(35),
	}
	for k, v := range table {
		p.put(k, v)
	}
	paths := p.walk()
	if len(paths) != len(table) {
		t.Errorf("expected %d got %d", len(table), len(paths))
	}
	p.printPaths()
}

func TestTriewalknormal(t *testing.T) {
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
		"delete":             dumpHandler(16),
	}
	for k, v := range table {
		p.put(k, v)
	}
	paths := p.walk()
	if len(paths) != len(table) {
		t.Errorf("expected %d got %d", len(table), len(paths))
	}
	p.printPaths()
}
