package http

import (
	"context"

	"github.com/valyala/fasthttp"
)

// Request is used to retrieve data from an HTTP request
type Request interface {
	// Path returns the path of the current URL
	Path() []byte

	// Method returns the HTTP method
	Method() []byte

	// IsJSON return weather request accepts application/json
	IsJSON() bool

	// URI returns the raw URI
	URI() *fasthttp.URI

	// Header return a header value by name. If the header is not found
	// an empty string will be returned.
	Header(name string) []byte

	// Host returns the host of the request.
	Host() []byte

	// Param grabs route param by name
	Param(name string) string

	// Query grabs input from the query string by name
	Query(name string) []byte

	// QueryMulti grabs multiple input from the query string by name
	QueryMulti(name string) [][]byte

	// Data unmarshals request body to dst
	Data(dst interface{}) error

	// Post grabs input from the post data by name
	Post(name string) []byte

	// PostMulti grabs multiple input from the post data by name
	PostMulti(name string) [][]byte

	// Cookie grabs input from cookies by name
	Cookie(name string) []byte

	// Context returns the context.Context of the current request
	Context() context.Context

	// WithContext returns a shallow copy of the request with a new context
	WithContext(ctx context.Context) Request

	// Raw returns the fasthttp.RequestCtx of the current request
	Raw() *fasthttp.RequestCtx
}

// Response is used to send data to the client
type Response interface {
	// Cookie sets an HTTP cookie on the response
	// See also `fasthttp.AcquireCookie`
	Cookie(cookie *fasthttp.Cookie) Response

	// Status sets the HTTP status code of the response. This can only be called once.
	Status(status int) Response

	// Header adds an HTTP header to the response
	Header(name, value string) Response

	// Data responds with data provided
	//
	// Most types will converted to a string representation except structs,
	// arrays and maps which will be serialized to JSON.
	Data(data interface{}) Result

	// Error sends the default 500 response
	Error(error) Result

	// Redirect redirects the client to a URL
	Redirect(uri string, code int) Result

	// End ends the response chain
	End() Result
}

// Result is used to finish a request
type Result interface {
	Data(data interface{}) Result
	Error(error) Result
	Redirect(uri string, code int) Result
}

type Handler func(req Request, res Response) Result

// Middleware is an interface for adding middleware to a Router instance
type Middleware func(req Request, res Response, next Handler) Result

type Router interface {
	Routable

	Handler() fasthttp.RequestHandler
}
