package hermes

type Routable interface {
	Delete(path string, handler Handler)
	Get(path string, handler Handler)
	Head(path string, handler Handler)
	Options(path string, handler Handler)
	Patch(path string, handler Handler)
	Post(path string, handler Handler)
	Put(path string, handler Handler)

	Prefix(path string) Routable
	Group(func(Routable))

	Use(...Middleware)
	With(...Middleware) Routable
}
