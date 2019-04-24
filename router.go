package http

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

type Router struct {
	route
	children map[string]*node
	NotFound Handler
}

type route struct {
	prefix      string
	router      *Router
	middlewares []Middleware
}

func NewRouter() *Router {
	r := &Router{
		children: make(map[string]*node),
	}
	r.router = r
	return r
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

func Split(source []byte, dest [][]byte) [][]byte {
	lSource := len(source)
	s := 0
	for i := 0; i < lSource; i++ {
		if source[i] == '/' {
			if i != s {
				dest = append(dest, source[s:i])
			}
			s = i + 1
		} else if i+1 == lSource {
			if i != s {
				dest = append(dest, source[s:i+1])
			}
		}
	}
	return dest
}

var routerHandlerSep = []byte{'/'}

var pathPool = sync.Pool{
	New: func() interface{} {
		return make([][]byte, 0, 255)
	},
}

func (router *Router) Handler(fCtx *fasthttp.RequestCtx) {
	method := string(fCtx.Method())
	node, ok := router.children[method]
	ctx := &Context{
		Ctx:      fCtx,
		Request:  &fCtx.Request,
		Response: &fCtx.Response,
	}
	if ok {
		path := pathPool.Get().([][]byte)
		path = Split(ctx.Request.URI().Path(), path)
		defer func() {
			path = path[0:0]
			pathPool.Put(path)
		}()
		path = bytes.Split(ctx.Request.URI().Path()[1:], routerHandlerSep)
		if len(path) == 1 && len(path[0]) == 0 {
			if node.handler != nil {
				node.handler(ctx)
				return
			}
			if router.NotFound != nil {
				router.NotFound(ctx)
			}
			return
		}
		found, node, values := node.Matches(path, nil)
		if found {
			for i, v := range values {
				ctx.Ctx.SetUserValue(node.names[i], string(v))
			}
			if len(node.middlewares) > 0 {
				middlewareIdx := 0
				var next Handler
				next = func(ctx *Context) {
					middlewareIdx++
					if middlewareIdx < len(node.middlewares) {
						node.middlewares[middlewareIdx](ctx, next)
					} else {
						node.handler(ctx)
					}
				}
				node.middlewares[0](ctx, next)
			} else {
				node.handler(ctx)
			}
			return
		}
	}
	if router.NotFound != nil {
		router.NotFound(ctx)
	}
}

func (group *route) DELETE(path string, handler Handler) {
	group.handle("DELETE", path, handler)
}

func (group *route) GET(path string, handler Handler) {
	group.handle("GET", path, handler)
}

func (group *route) POST(path string, handler Handler) {
	group.handle("POST", path, handler)
}

func (group *route) PUT(path string, handler Handler) {
	group.handle("PUT", path, handler)
}

func (group *route) HEAD(path string, handler Handler) {
	group.handle("HEAD", path, handler)
}

func (group *route) OPTIONS(path string, handler Handler) {
	group.handle("OPTIONS", path, handler)
}

func (group *route) PATCH(path string, handler Handler) {
	group.handle("PATCH", path, handler)
}

func (group *route) Prefix(path string) Routable {
	return &route{
		prefix:      fmt.Sprintf("%s%s", group.prefix, path),
		router:      group.router,
		middlewares: group.middlewares,
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
