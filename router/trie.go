package router

import (
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
	pathSeperator = "/"
	paramNote     = ":"
	endChildKey   = ":END" + pathSeperator + ":"
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
	if path == "" {
		return nil
	}
	if p.children == nil {
		return nil
	}
	ks := strings.Split(path, pathSeperator)
	key := ks[0]
	child, ok := p.children[key]
	if ok {
		if len(ks) == 1 {
			return child.value
		}
		return child.get(strings.Join(ks[1:], pathSeperator))
	}
	return nil
}

// suffix '/' will be trimed
func (p *pathTrie) put(path string, value http.Handler) *pathTrie {
	path = strings.Trim(strings.TrimSpace(path), pathSeperator)
	if path == "" {
		return nil
	}
	if p.children == nil {
		p.children = make(map[string]*pathTrie)
	}
	node := newPathTrie()
	node.value = value
	ks := strings.Split(path, pathSeperator)
	key := ks[0]
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

func (p *pathTrie) walk(prefix string, n int) {
	if p == nil {
		fmt.Println(nil)
		return
	}
	if p.value == nil {
		// fmt.Print(nil)
	} else {
		fmt.Printf(" --> %#v", p.value)
	}
	if len(p.children) == 0 {
		return
	}
	for k, v := range p.children {
		if strings.HasPrefix(k, pathSeperator) {
			k = trimPathParam(k)
		}
		sub := fmt.Sprintf("%s/%s", prefix, k)
		if v.value != nil {
			fmt.Printf("\n")
			fmt.Print(sub)
		}
		v.walk(sub, n+1)
	}
}

func trimPathParam(k string) string {
	param := strings.TrimPrefix(k, pathSeperator)
	return paramNote + param
}

func makePathParamKey(param string) string {
	k := strings.TrimPrefix(param, paramNote)
	return pathSeperator + k
}
