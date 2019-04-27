package middlewares

import (
	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/lab259/http"
	"github.com/valyala/fasthttp"
)

func DefaultCorsMiddleware() http.Middleware {
	return wrapCorsMiddleware(cors.DefaultHandler())
}

func NewCorsMiddleware(options cors.Options) http.Middleware {
	return wrapCorsMiddleware(cors.NewCorsHandler(options))
}

func wrapCorsMiddleware(withCors *cors.CorsHandler) http.Middleware {
	return func(req http.Request, res http.Response, next http.Handler) http.Result {
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
