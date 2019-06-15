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

func DefaultRouter() Router {
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
				errResponse := acquireErrorResponse(StatusNotFound)
				defer releaseErrorResponse(errResponse)

				errResponse.SetParam("code", NotFoundErrorCode)
				errResponse.SetParam("message", NotFoundErrorMessage)
				return res.Status(errResponse.Status).Data(errResponse.Data)
			}

			return res.Status(StatusNotFound).Data(NotFoundErrorMessage)
		}
	}

	if config.MethodNotAllowed == nil {
		r.methodNotAllowed = func(req Request, res Response) Result {
			if req.WantsJSON() {
				errResponse := acquireErrorResponse(StatusMethodNotAllowed)
				defer releaseErrorResponse(errResponse)

				errResponse.SetParam("code", MethodNotAllowedErrorCode)
				errResponse.SetParam("message", MethodNotAllowedErrorMessage)
				return res.Status(errResponse.Status).Data(errResponse.Data)
			}

			return res.Status(StatusMethodNotAllowed).Data(MethodNotAllowedErrorMessage)
		}
	}

	r.defaultOptions = func(req Request, res Response) Result {
		return res.End()
	}

	r.router = r
	return r
}

var routerHandlerSep = []byte{'/'}

type tokensDescriptor struct {
	m [][]byte
	n int
}

var tokensDescriptorPool = sync.Pool{
	New: func() interface{} {
		return &tokensDescriptor{
			m: make([][]byte, 0, 10),
			n: 0,
		}
	},
}

func acquireTokensDescriptor() *tokensDescriptor {
	return tokensDescriptorPool.Get().(*tokensDescriptor)
}

func releaseTokensDescriptor(path *tokensDescriptor) {
	path.n = 0
	path.m = path.m[:0]
	tokensDescriptorPool.Put(path)
}

func split(source []byte, dest *tokensDescriptor) {
	lSource := len(source)
	s := 0
	for i := 0; i < lSource; i++ {
		if source[i] == '/' {
			if i != s {
				dest.m = append(dest.m, source[s:i])
				dest.n++
			}
			s = i + 1
		} else if i+1 == lSource {
			if i != s {
				dest.m = append(dest.m, source[s:i+1])
				dest.n++
			}
		}
	}
}

func (router *router) findHandler(root *node, path *tokensDescriptor, values *tokensDescriptor) (bool, *node) {
	if path.n == 0 {
		if root.handler != nil {
			return true, root
		}
	}
	return root.Matches(0, path, values)
}

func (router *router) releaseResources(req *request, res *response, path *tokensDescriptor, values *tokensDescriptor) {
	releaseRequest(req)
	releaseResponse(res)
	releaseTokensDescriptor(path)
	releaseTokensDescriptor(values)
}

func (router *router) Handler() fasthttp.RequestHandler {
	return func(fCtx *fasthttp.RequestCtx) {
		req := acquireRequest(context.Background(), fCtx)
		res := acquireResponse(fCtx)
		values := acquireTokensDescriptor()
		path := acquireTokensDescriptor()
		defer router.releaseResources(req, res, path, values)

		// split request path into tokenDescriptor
		split(req.Path(), path)

		method := string(req.Method())
		if root, ok := router.children[method]; ok {
			if found, node := router.findHandler(root, path, values); found {
				req.params = values.m[:]
				req.validParams = node.names[:]
				node.handler(req, res).End()
				return
			}
		}

		if method == "OPTIONS" {
			// handle OPTIONS requests
			if allow := router.allowed(req.Path(), method, path); len(allow) > 0 {
				res.Header("Allow", allow)
				router.callHandler(req, res, router.middlewares, router.defaultOptions).End()
				return
			}
		} else {
			// handle 405
			if allow := router.allowed(req.Path(), method, path); len(allow) > 0 {
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

func (r *router) allowed(reqPath []byte, reqMethod string, path *tokensDescriptor) (allow string) {
	if bytes.Equal(reqPath, optionsServerWide) || bytes.Equal(reqPath, optionsSlashServerWide) { // server-wide
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

			if found, _ := r.findHandler(r.children[method], path, nil); found {
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
