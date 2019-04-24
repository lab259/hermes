package http

type Routable interface {
	handle(method, path string, handler Handler)

	DELETE(path string, handler Handler)
	GET(path string, handler Handler)
	HEAD(path string, handler Handler)
	OPTIONS(path string, handler Handler)
	PATCH(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)

	Prefix(path string) Routable
	Group(func(Routable))

	Use(...Middleware)
	With(...Middleware) Routable
}
