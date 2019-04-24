package http

// Middleware is the descriptor of a middleware that wraps a handler.
//
// A middleware implementation will receive a `*Context` and a `next` as a
// representation of the next middleware (or the end handler itself) that should
// be called if the endpoint execution should carry on.
//
// If you need to avoid the endpoint execution to continue, you should not to
// call the `next` handler.
// type Middleware func(ctx *Context, next Handler)
