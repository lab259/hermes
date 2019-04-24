package http

import (
	"bytes"
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

type router struct {
	route
	children map[string]*node
	notFound Handler
}

func NewRouter(notFound Handler) Router {
	r := &router{
		children: make(map[string]*node),
		notFound: notFound,
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

		method := string(fCtx.Method())
		node, ok := router.children[method]

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
				if router.notFound != nil {
					router.notFound(req, res)
				}
				return
			}
			found, node, values := node.Matches(path, nil)
			if found {
				for i, v := range values {
					fCtx.SetUserValue(node.names[i], string(v))
				}
				if len(node.middlewares) > 0 {
					middlewareIdx := 0
					var next Handler
					next = func(req2 Request, res2 Response) Result {
						middlewareIdx++
						if middlewareIdx < len(node.middlewares) {
							return node.middlewares[middlewareIdx](req2, res2, next)
						}

						return node.handler(req2, res2)
					}
					node.middlewares[0](req, res, next)
				} else {
					node.handler(req, res)
				}
				return
			}
		}
		if router.notFound != nil {
			router.notFound(req, res)
		}
	}
}
