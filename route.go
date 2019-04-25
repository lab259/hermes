package http

import "fmt"

type route struct {
	prefix      string
	router      *router
	middlewares []Middleware
}

func (r *route) path(subpath string) string {
	if len(subpath) > 0 && subpath[0] == '/' {
		subpath = subpath[1:]
	}
	if subpath == "" {
		return r.prefix
	}
	if r.prefix == "" {
		return subpath
	}
	return fmt.Sprintf("%s/%s", r.prefix, subpath)
}

func (r *route) handle(method, subpath string, handler Handler) {
	root, ok := r.router.children[method]
	if !ok {
		root = newNode()
		r.router.children[method] = root
	}
	path := r.path(subpath)
	root.Add(path, handler, nil, r.middlewares)
}

func (r *route) Delete(path string, handler Handler) {
	r.handle("DELETE", path, handler)
}

func (r *route) Get(path string, handler Handler) {
	r.handle("GET", path, handler)
}

func (r *route) Post(path string, handler Handler) {
	r.handle("POST", path, handler)
}

func (r *route) Put(path string, handler Handler) {
	r.handle("PUT", path, handler)
}

func (r *route) Head(path string, handler Handler) {
	r.handle("HEAD", path, handler)
}

func (r *route) Options(path string, handler Handler) {
	r.handle("OPTIONS", path, handler)
}

func (r *route) Patch(path string, handler Handler) {
	r.handle("PATCH", path, handler)
}

func (r *route) Prefix(path string) Routable {
	return &route{
		prefix:      r.path(path),
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
