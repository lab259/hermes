package http

import (
	"bytes"
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

type RouterConfig struct {
	NotFoundHandler Handler
}

type router struct {
	route
	children map[string]*node
	notFound Handler
}

func NewRouter(config RouterConfig) Router {
	r := &router{
		children: make(map[string]*node),
		notFound: config.NotFoundHandler,
	}
	r.router = r
	return r
}

func split(source []byte, dest [][]byte) [][]byte {
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

func (router *router) Handler() fasthttp.RequestHandler {
	return func(fCtx *fasthttp.RequestCtx) {
		req := acquireRequest(context.Background(), fCtx)
		defer releaseRequest(req)

		res := acquireResponse(fCtx)
		defer releaseResponse(res)

		node, ok := router.children[string(req.Method())]
		if ok {
			path := pathPool.Get().([][]byte)
			path = split(req.Path(), path)
			defer func() {
				path = path[0:0]
				pathPool.Put(path)
			}()
			path = bytes.Split(req.Path()[1:], routerHandlerSep)
			if len(path) == 1 && len(path[0]) == 0 {
				if node.handler != nil {
					node.handler(req, res)
					return
				}
				router.callNotFound(req, res)
				return
			}
			found, node, values := node.Matches(path, nil)
			if found {
				for i, v := range values {
					fCtx.SetUserValue(node.names[i], string(v))
				}
				node.handler(req, res)
				return
			}
		}
		router.callNotFound(req, res)
	}
}

func (router *router) callHandler(req Request, res Response, middlewares []Middleware, handler Handler) {
	if len(middlewares) > 0 {
		middlewareIdx := 0
		var next Handler
		next = func(req2 Request, res2 Response) Result {
			middlewareIdx++
			if middlewareIdx < len(middlewares) {
				return middlewares[middlewareIdx](req2, res2, next)
			}

			return handler(req2, res2)
		}
		middlewares[0](req, res, next)
	} else {
		handler(req, res)
	}
}

func (router *router) callNotFound(req Request, res Response) {
	if router.notFound != nil {
		router.callHandler(req, res, router.middlewares, router.notFound)
	}
}
