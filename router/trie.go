package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type pathRegs string

const (
	paramsCtxKey   pathRegs = "path_param_names"
	paramsCtxValue pathRegs = "path_param_values"
)
const (
	pathSeperator  = "/"
	paramNote      = ":"
	regexpChildKey = ":REG" + pathSeperator + ":"
	endChildKey    = ":END" + pathSeperator + ":"
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

// get return immediately when match endChildKey, but do NOT
// ignore plain child(include regex child) on the same level
func (p *pathTrie) get(path string) http.Handler {
	parts := strings.Split(strings.TrimSpace(path), pathSeperator)
	node := p
	paramValues := []string{}
	wrapHandler := func(h http.Handler, vs []string) http.Handler {
		if len(vs) > 0 && h != nil {
			fn := func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), paramsCtxValue, paramValues)
				r = r.WithContext(ctx)
				h.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		}
		return h
	}
	for _, v := range parts {
		child, ok := node.children[v]
		if !ok {
			if s, ok := node.children[regexpChildKey]; ok {
				paramValues = append(paramValues, v)
				node = s
				continue
			}
			if e, ok := node.children[endChildKey]; ok {
				return wrapHandler(e.value, paramValues)
			}
			return nil
		}
		node = child
	}
	return wrapHandler(node.value, paramValues)
}

// suffix '/' counts: path '/a' is diffent from path '/a/'
func (p *pathTrie) put(path string, value http.Handler) *pathTrie {
	parts := strings.Split(strings.TrimSpace(path), pathSeperator)
	node := p
	regs := []string{}
	for _, v := range parts {
		child, ok := node.children[v]
		if strings.HasPrefix(v, paramNote) {
			regs = append(regs, strings.TrimPrefix(v, paramNote))
			v = regexpChildKey
		}
		if !ok {
			if special, o := node.children[regexpChildKey]; o {
				child = special
			} else {
				child = &pathTrie{
					value:    nil,
					children: make(map[string]*pathTrie),
				}
			}
			node.children[v] = child
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
		return node
	}
	node.value = value
	return node
}

// putEnd stop the get the the endChildKey
// call put then append a endChild to node
// useful for http.Fileserver
func (p *pathTrie) putEnd(path string, value http.Handler) {
	node := p.put(path, value)
	node.children[endChildKey] = &pathTrie{
		value:    value,
		children: make(map[string]*pathTrie),
	}
}

func (p *pathTrie) walk(n int) {
	if p == nil {
		fmt.Println(nil)
		return
	}
	fmt.Printf("handler:%q,%d children.", p.value, len(p.children))
	if len(p.children) == 0 {
		return
	}
	for k, v := range p.children {
		fmt.Printf("\n" + strings.Repeat("\t", n))
		fmt.Printf("%q->", k)
		v.walk(n + 1)
	}
}
