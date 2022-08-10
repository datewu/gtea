package router

import (
	"net/http"
	"strings"
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
	parts := strings.Split(path, "/")
	node := p
	for _, v := range parts {
		child, ok := node.children[v]
		if !ok {
			return nil
		} else {
			node = child
		}
	}
	return node.value
}

// suffix '/' counts: path '/a' is diffent from path '/a/'
func (p *pathTrie) put(path string, value http.Handler) {
	parts := strings.Split(path, "/")
	node := p
	for _, v := range parts {
		// not skip '/' suffix
		// if v == "" {
		// 	continue
		// }
		child, ok := node.children[v]
		if !ok {
			child = &pathTrie{
				value:    nil,
				children: make(map[string]*pathTrie),
			}
			node.children[v] = child
		}
		node = child
	}
	if node.value != nil {
		panic("trie node value conflict")
	}
	node.value = value
}

// func (p *pathTrie) delete(key string) bool {
// 	return true // node (internal or not) existed and its value was nil'd
// }
