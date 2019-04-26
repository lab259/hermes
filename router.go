package http

import (
	"bytes"
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

type RouterConfig struct {
	NotFound         Handler
	MethodNotAllowed Handler
}

type router struct {
	route
	children         map[string]*node
	notFound         Handler
	methodNotAllowed Handler
	defaultOptions   Handler
}

func NewDefaultRouter() Router {
	return NewRouter(RouterConfig{})
}

func NewRouter(config RouterConfig) Router {
	r := &router{
		children:         make(map[string]*node),
		notFound:         config.NotFound,
		methodNotAllowed: config.MethodNotAllowed,
	}

	if config.NotFound == nil {
		r.notFound = func(req Request, res Response) Result {
			if req.WantsJSON() {
				errResponse := acquireErrorResponse(fasthttp.StatusNotFound)
				defer releaseErrorResponse(errResponse)

				errResponse.SetParam("code", NotFoundErrorCode)
				errResponse.SetParam("message", NotFoundErrorMessage)
				return res.Status(errResponse.Status).Data(errResponse.Data)
			}

			return res.Status(fasthttp.StatusNotFound).Data(NotFoundErrorMessage)
		}
	}

	if config.MethodNotAllowed == nil {
		r.methodNotAllowed = func(req Request, res Response) Result {
			if req.WantsJSON() {
				errResponse := acquireErrorResponse(fasthttp.StatusMethodNotAllowed)
				defer releaseErrorResponse(errResponse)

				errResponse.SetParam("code", MethodNotAllowedErrorCode)
				errResponse.SetParam("message", MethodNotAllowedErrorMessage)
				return res.Status(errResponse.Status).Data(errResponse.Data)
			}

			return res.Status(fasthttp.StatusMethodNotAllowed).Data(MethodNotAllowedErrorMessage)
		}
	}

	r.defaultOptions = func(req Request, res Response) Result {
		return res.End()
	}

	r.router = r
	return r
}

var routerHandlerSep = []byte{'/'}

type pathDescriptor struct {
	m map[int][]byte
	n int
}

var pathPool = sync.Pool{
	New: func() interface{} {
		return &pathDescriptor{
			m: make(map[int][]byte, 10),
			n: 0,
		}
	},
}

func acquirePath() *pathDescriptor {
	return pathPool.Get().(*pathDescriptor)
}

func releasePath(path *pathDescriptor) {
	path.n = 0
	pathPool.Put(path)
}

func split(source []byte, dest *pathDescriptor) {
	lSource := len(source)
	s := 0
	for i := 0; i < lSource; i++ {
		if source[i] == '/' {
			if i != s {
				dest.m[dest.n] = source[s:i]
				dest.n++
			}
			s = i + 1
		} else if i+1 == lSource {
			if i != s {
				dest.m[dest.n] = source[s : i+1]
				dest.n++
			}
		}
	}
}

func (router *router) findHandler(root *node, reqPath []byte) (bool, *node, [][]byte) {
	path := acquirePath()
	defer releasePath(path)

	split(reqPath, path)
	if path.n == 0 {
		if root.handler != nil {
			return true, root, nil
		}
	}
	return root.Matches(0, path, nil)
}

func (router *router) Handler() fasthttp.RequestHandler {
	return func(fCtx *fasthttp.RequestCtx) {
		req := acquireRequest(context.Background(), fCtx)
		defer releaseRequest(req)

		res := acquireResponse(fCtx)
		defer releaseResponse(res)

		method := string(req.Method())
		if root, ok := router.children[method]; ok {
			if found, node, values := router.findHandler(root, req.Path()); found {
				for i, v := range values {
					fCtx.SetUserValue(node.names[i], string(v))
				}
				node.handler(req, res).End()
				return
			}
		}

		if method == "OPTIONS" {
			// handle OPTIONS requests
			if allow := router.allowed(req.Path(), method); len(allow) > 0 {
				res.Header("Allow", allow)
				router.callHandler(req, res, router.middlewares, router.defaultOptions).End()
				return
			}
		} else {
			// handle 405
			if allow := router.allowed(req.Path(), method); len(allow) > 0 {
				res.Header("Allow", allow)
				router.callHandler(req, res, router.middlewares, router.methodNotAllowed).End()
				return
			}
		}

		router.callHandler(req, res, router.middlewares, router.notFound).End()
	}
}

func (router *router) callHandler(req Request, res Response, middlewares []Middleware, handler Handler) Result {
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
		return middlewares[0](req, res, next)
	}
	return handler(req, res)
}

var (
	optionsServerWide      = []byte("*")
	optionsSlashServerWide = []byte("/*")
)

func (r *router) allowed(path []byte, reqMethod string) (allow string) {
	if bytes.Equal(path, optionsServerWide) || bytes.Equal(path, optionsSlashServerWide) { // server-wide
		for method := range r.children {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.children {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			if found, _, _ := r.findHandler(r.children[method], path); found {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}
