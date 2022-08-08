package router

import "testing"

func TestPathReg(t *testing.T) {
	result := map[string][]string{
		"/":                     nil,
		"/a/ba/c":               nil,
		"/a/:name":              {"name"},
		"/a/:name/b/:age":       {"name", "age"},
		"/a/:name_after/:title": {"name_after", "title"},
		"/:a/:b/:c/d":           {"a", "b", "c"},
	}
	for k, v := range result {
		r := findPathParam(k)
		if v == nil {
			if containerPathParam(k) {
				t.Error("should have no path params")
			}
			if r != nil {
				t.Error("should have no path params", "got:", r)
				continue
			}
		}
		if len(r) != len(v) {
			t.Error("path params length not matched:", v, r)
			continue
		}
		for i := range r {
			if r[i] != v[i] {
				t.Error("got wrong params:", v, r)
			}
		}
	}

}
