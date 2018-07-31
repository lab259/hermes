package http

type Routable interface {
	DELETE(path string, handler Handler, middlewares ...Middleware)
	GET(path string, handler Handler, middlewares ...Middleware)
	HEAD(path string, handler Handler, middlewares ...Middleware)
	OPTIONS(path string, handler Handler, middlewares ...Middleware)
	PATCH(path string, handler Handler, middlewares ...Middleware)
	POST(path string, handler Handler, middlewares ...Middleware)
	PUT(path string, handler Handler, middlewares ...Middleware)

	Group(path string, middlewares ...Middleware) Routable
}
