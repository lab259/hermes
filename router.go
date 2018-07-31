package http

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"sync"
	"bytes"
)

type Router struct {
	children map[string]*node
	NotFound Handler
}

func NewRouter() *Router {
	return &Router{
		children: make(map[string]*node),
	}
}

func (router *Router) handle(method, path string, handler Handler, middlewares ...Middleware) {
	root, ok := router.children[method]
	if !ok {
		root = newNode()
		router.children[method] = root
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	root.Add(path, handler, nil, middlewares)
}

func (router *Router) DELETE(path string, handler Handler, middlewares ...Middleware) {
	router.handle("DELETE", path, handler, middlewares...)
}

func (router *Router) GET(path string, handler Handler, middlewares ...Middleware) {
	router.handle("GET", path, handler, middlewares...)
}

func (router *Router) POST(path string, handler Handler, middlewares ...Middleware) {
	router.handle("POST", path, handler, middlewares...)
}

func (router *Router) PUT(path string, handler Handler, middlewares ...Middleware) {
	router.handle("PUT", path, handler, middlewares...)
}

func (router *Router) HEAD(path string, handler Handler, middlewares ...Middleware) {
	router.handle("HEAD", path, handler, middlewares...)
}

func (router *Router) OPTIONS(path string, handler Handler, middlewares ...Middleware) {
	router.handle("OPTIONS", path, handler, middlewares...)
}

func (router *Router) PATCH(path string, handler Handler, middlewares ...Middleware) {
	router.handle("PATCH", path, handler, middlewares...)
}

func (router *Router) Group(path string, middlewares ... Middleware) Routable {
	return &routerGroup{
		prefix:      path,
		router:      router,
		middlewares: middlewares,
	}
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

type routerGroup struct {
	prefix      string
	router      Routable
	middlewares []Middleware
}

func (group *routerGroup) DELETE(path string, handler Handler, middlewares ...Middleware) {
	group.router.DELETE(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) GET(path string, handler Handler, middlewares ...Middleware) {
	group.router.GET(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) POST(path string, handler Handler, middlewares ...Middleware) {
	group.router.POST(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) PUT(path string, handler Handler, middlewares ...Middleware) {
	group.router.PUT(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) HEAD(path string, handler Handler, middlewares ...Middleware) {
	group.router.HEAD(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) OPTIONS(path string, handler Handler, middlewares ...Middleware) {
	group.router.OPTIONS(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) PATCH(path string, handler Handler, middlewares ...Middleware) {
	group.router.PATCH(fmt.Sprintf("%s%s", group.prefix, path), handler, append(group.middlewares, middlewares...)...)
}

func (group *routerGroup) Group(path string, middlewares ... Middleware) Routable {
	return &routerGroup{
		prefix:      path,
		router:      group,
		middlewares: middlewares,
	}
}
