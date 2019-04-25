package http

import (
	"bytes"
	"fmt"
)

type node struct {
	wildcard *node
	children map[string]*node
	handler  Handler
	names    []string
}

func newNode() *node {
	return &node{
		children: make(map[string]*node),
	}
}

func (n *node) Add(path string, handler Handler, names []string, middlewares []Middleware) {
	// Split path into chunks between `/`
	pathBytes := bytes.Split([]byte(path), []byte{'/'})
	lpath := len(pathBytes)
	parent := n
	// Go though all chars
	for i := 0; i < lpath; i++ {
		token := pathBytes[i]
		if len(token) > 0 {
			// If this token is a wildcard
			if token[0] == ':' {
				name := string(token[1:])
				node := parent.wildcard
				nodeCreated := false
				if node == nil {
					node = newNode()
					parent.wildcard = node
					nodeCreated = true
				}
				// If names is not defined yet
				if names == nil {
					names = make([]string, 0, 2) // Starts with a capacity of 2
				}
				names = append(names, name)
				parent = node
				if i+1 >= lpath {
					if !nodeCreated && node.handler != nil {
						// Two wildcard created with the same stuff
						panic(fmt.Sprintf("conflict adding '%s'", path))
					}
					// Initialize stuff
					node.names = names
					node.handler = newHandler(handler, middlewares)
				}
				continue
			} else {
				spath := string(pathBytes[i])
				node, ok := parent.children[spath]
				// If the child does not exists, it will create a new one.
				if !ok {
					node = newNode()
					parent.children[spath] = node
				}

				// If this is not the end
				if i+1 < lpath {
					parent = node
				} else {
					// This is the end of the path
					if ok && node.handler != nil {
						panic(fmt.Sprintf("conflict adding '%s'", path))
					}
					node.names = names
					node.handler = newHandler(handler, middlewares)
					return
				}
			}
		} else if i+1 < lpath {
			// Cannot deal with empty tokens
			panic("empty token")
		} else {
			// Just set the node info
			n.names = names
			n.handler = newHandler(handler, middlewares)
		}
	}
}

func newHandler(h Handler, m []Middleware) Handler {
	sagas := make([]Handler, len(m)+1)

	sagas[len(m)] = h

	idx := len(m) - 1
	for idx > -1 {
		i := idx
		sagas[i] = func(nextReq Request, nextRes Response) Result {
			return m[i](nextReq, nextRes, sagas[i+1])
		}
		idx--
	}

	return sagas[0]
}

func (n *node) Matches(path [][]byte, values [][]byte) (bool, *node, [][]byte) {
	lpath := len(path)
	for i := 0; i < lpath; i++ {
		token := string(path[i])
		node, ok := n.children[token]
		if ok {
			if i+1 < lpath {
				return node.Matches(path[i+1:], values)
			} else if node.handler == nil {
				return false, nil, nil
			} else {
				return true, node, values
			}
		} else if n.wildcard != nil {
			if values == nil {
				values = [][]byte{path[i]}
			} else {
				values = append(values, path[i])
			}
			if i+1 < lpath {
				return n.wildcard.Matches(path[i+1:], values)
			} else if n.wildcard.handler == nil {
				return false, nil, nil
			} else {
				return true, n.wildcard, values
			}
		} else {
			return false, nil, nil
		}
	}
	return false, nil, nil
}
