package middlewares

import (
	"github.com/lab259/cors"
	"github.com/lab259/hermes"
	"github.com/valyala/fasthttp"
)

func DefaultCorsMiddleware() hermes.Middleware {
	return wrapCorsMiddleware(cors.Default())
}

func NewCorsMiddleware(options cors.Options) hermes.Middleware {
	return wrapCorsMiddleware(cors.New(options))
}

func wrapCorsMiddleware(withCors *cors.Cors) hermes.Middleware {
	return func(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
		canContinue := false
		withCors.Handler(func(ctx *fasthttp.RequestCtx) {
			canContinue = true
		})(req.Raw())
		if canContinue {
			return next(req, res)
		}
		return res.End()
	}
}
