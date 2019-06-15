package middlewares

import (
	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/lab259/hermes"
	"github.com/valyala/fasthttp"
)

func DefaultCorsMiddleware() hermes.Middleware {
	return wrapCorsMiddleware(cors.DefaultHandler())
}

func NewCorsMiddleware(options cors.Options) hermes.Middleware {
	return wrapCorsMiddleware(cors.NewCorsHandler(options))
}

func wrapCorsMiddleware(withCors *cors.CorsHandler) hermes.Middleware {
	return func(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
		canContinue := false
		withCors.CorsMiddleware(func(ctx *fasthttp.RequestCtx) {
			canContinue = true
		})(req.Raw())
		if canContinue {
			return next(req, res)
		}
		return res.End()
	}
}
