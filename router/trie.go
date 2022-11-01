package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/datewu/gtea/handler"
)

const (
	pathSeperator = "/"
	paramNote     = ":"
	regKey        = paramNote + pathSeperator + "REG"
	// end of level, no descendants/children anymore
	endChildKey = regKey + "EOL"
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
	path = strings.Trim(strings.TrimSpace(path), pathSeperator)
	if path == "" || p.children == nil {
		return p.value
	}
	vs := strings.Split(path, pathSeperator)
	value := vs[0]
	child, ok := p.children[value]
	if ok {
		if len(vs) == 1 {
			return child.value
		}
		return child.get(strings.Join(vs[1:], pathSeperator))
	}
	regChild, ok := p.children[regKey]
	if ok {
		if len(vs) == 1 {
			return setParamValue(regChild.value, value)
		}
		return setParamValue(regChild.get(strings.Join(vs[1:], pathSeperator)), value)
	}
	endChild, ok := p.children[endChildKey]
	if ok {
		return endChild.value
	}
	return nil
}

func insertCtxValue(v http.Handler, key handler.PathRegs, value string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var state []string
		vs := r.Context().Value(key)
		if vs == nil {
			state = []string{value}
		} else {
			data, ok := vs.([]string)
			if !ok {
				panic("should be []string in " + key)
			}
			tmp := make([]string, len(data)+1)
			tmp[0] = value
			for i, v := range data {
				tmp[i+1] = v
			}
			state = tmp
		}
		ctx := context.WithValue(r.Context(), key, state)
		r = r.WithContext(ctx)
		v.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func addCtxValue(v http.Handler, key handler.PathRegs, value string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var state []string
		vs := r.Context().Value(key)
		if vs == nil {
			state = []string{value}
		} else {
			data, ok := vs.([]string)
			if !ok {
				panic("should be []string in " + key)
			}
			data = append(data, value)
			state = data
		}
		ctx := context.WithValue(r.Context(), key, state)
		r = r.WithContext(ctx)
		v.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setParamValue(v http.Handler, value string) http.Handler {
	return addCtxValue(v, handler.ParamsCtxValue, value)
}

func addParamName(v http.Handler, name string) http.Handler {
	return insertCtxValue(v, handler.ParamsCtxKey, name)
}

// suffix '/' will be trimed
func (p *pathTrie) put(path string, value http.Handler) *pathTrie {
	path = strings.Trim(strings.TrimSpace(path), pathSeperator)
	if path == "" {
		p.value = value
		return p
	}
	if p.children == nil {
		panic("no children, maybe on EOL no descendants")
	}
	ks := strings.Split(path, pathSeperator)
	key := ks[0]
	if strings.HasPrefix(key, paramNote) {
		value = addParamName(value, key[1:])
		key = regKey
	}
	node := newPathTrie()
	node.value = value
	child, ok := p.children[key]
	if ok {
		if len(ks) == 1 {
			if child.value != nil {
				panic("node conflict")
			} else {
				child.value = value
				return child
			}
		}
		return child.put(strings.Join(ks[1:], pathSeperator), value)
	}
	if len(ks) == 1 {
		p.children[key] = node
		return node
	}
	p.children[key] = newPathTrie()
	return p.children[key].put(strings.Join(ks[1:], pathSeperator), value)
}

// putEnd end trie at that level
// stop look up children pathTrie
// useful for http.Fileserver wild path
func (p *pathTrie) putEnd(path string, value http.Handler) {
	node := p.put(path, value)
	node.children[endChildKey] = &pathTrie{
		value: value,
		//	children: make(map[string]*pathTrie),
	}
}

func (p *pathTrie) printPaths() {
	paths := p.walk()
	fmt.Printf("total %d paths: \nDetail:\n", len(paths))
	fmt.Println(strings.Join(paths, "\n"))
	fmt.Println("=======")
}

func (p *pathTrie) walk() []string {
	if p == nil {
		return nil
	}
	if p.children == nil {
		return nil
	}
	var result []string
	for line, child := range p.children {
		if child.value != nil {
			meat := fmt.Sprintf("%s --> %v", line, child.value)
			result = append(result, meat)
		}
		next := child.walk()
		for _, v := range next {
			result = append(result, line+pathSeperator+v)
		}
		if next == nil && child.value == nil {
			result = append(result, line)
		}
	}
	return result
}
