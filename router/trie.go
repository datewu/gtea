package router

import (
	"context"
	"net/http"
	"strings"
)

type pathRegs string
type pathName string

const paramsCtxKey pathRegs = "path_param_names"
const (
	pathSeperator   = "/"
	paramNote       = ":"
	specialChildKey = "  "
)

func newPathTrie() *pathTrie {
	root := &pathTrie{
		children: make(map[string]*pathTrie),
	}
	return root
}

// suffix '/' counts: path '/a' is diffent from path '/a/'
type pathTrie struct {
	value    http.Handler
	children map[string]*pathTrie
}

func (p *pathTrie) get(path string) http.Handler {
	parts := strings.Split(strings.TrimSpace(path), pathSeperator)
	node := p
	paramValues := []string{}
	for _, v := range parts {
		child, ok := node.children[v]
		if !ok {
			if s, ok := node.children[specialChildKey]; ok {
				paramValues = append(paramValues, v)
				node = s
				continue
			}
			return nil
		}
		node = child
	}
	if len(paramValues) > 0 {
		if node.value == nil {
			return nil
		}
		wrapH := func(w http.ResponseWriter, r *http.Request) {
			keys := r.Context().Value(paramsCtxKey)
			ks, ok := keys.([]string)
			if !ok || len(paramValues) != len(ks) {
				return
			}
			ctx := r.Context()
			for i := range ks {
				ctx = context.WithValue(ctx, pathName(ks[i]), paramValues[i])
			}
			r = r.WithContext(ctx)
			node.value.ServeHTTP(w, r)
		}
		return http.HandlerFunc(wrapH)
	}
	return node.value
}

// suffix '/' counts: path '/a' is diffent from path '/a/'
func (p *pathTrie) put(path string, value http.Handler) {
	parts := strings.Split(strings.TrimSpace(path), pathSeperator)
	node := p
	regs := []string{}
	for _, v := range parts {
		child, ok := node.children[v]
		if !ok {
			child = &pathTrie{
				value:    nil,
				children: make(map[string]*pathTrie),
			}
			if strings.HasPrefix(v, paramNote) {
				regs = append(regs, strings.TrimPrefix(v, ":"))
				node.children[specialChildKey] = child
			} else {
				node.children[v] = child
			}
		}
		node = child
	}
	if node.value != nil {
		panic("trie node value conflict")
	}
	if len(regs) > 0 {
		wrapH := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), paramsCtxKey, regs)
			r = r.WithContext(ctx)
			value.ServeHTTP(w, r)
		}
		node.value = http.HandlerFunc(wrapH)
		return
	}
	node.value = value
}

// func (p *pathTrie) delete(key string) bool {
// 	return true // node (internal or not) existed and its value was nil'd
// }
