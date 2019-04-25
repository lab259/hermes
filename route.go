package http

import "fmt"

type route struct {
	prefix      string
	router      *router
	middlewares []Middleware
}

func (r *route) handle(method, subpath string, handler Handler) {
	root, ok := r.router.children[method]
	if !ok {
		root = newNode()
		r.router.children[method] = root
	}
	path := fmt.Sprintf("%s%s", r.prefix, subpath)
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	root.Add(path, handler, nil, r.middlewares)
}

func (r *route) DELETE(path string, handler Handler) {
	r.handle("DELETE", path, handler)
}

func (r *route) GET(path string, handler Handler) {
	r.handle("GET", path, handler)
}

func (r *route) POST(path string, handler Handler) {
	r.handle("POST", path, handler)
}

func (r *route) PUT(path string, handler Handler) {
	r.handle("PUT", path, handler)
}

func (r *route) HEAD(path string, handler Handler) {
	r.handle("HEAD", path, handler)
}

func (r *route) OPTIONS(path string, handler Handler) {
	r.handle("OPTIONS", path, handler)
}

func (r *route) PATCH(path string, handler Handler) {
	r.handle("PATCH", path, handler)
}

func (r *route) Prefix(path string) Routable {
	return &route{
		prefix:      fmt.Sprintf("%s%s", r.prefix, path),
		router:      r.router,
		middlewares: r.middlewares,
	}
}

func (r *route) Group(h func(Routable)) {
	h(r)
}

func (r *route) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *route) With(middlewares ...Middleware) Routable {
	return &route{
		prefix:      r.prefix,
		router:      r.router,
		middlewares: append(r.middlewares, middlewares...),
	}
}