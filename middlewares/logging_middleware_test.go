package middlewares_test

import (
	"bytes"

	"github.com/lab259/rlog"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("Middlewares", func() {
	Describe("Logging Middleware", func() {
		It("should log requests", func() {
			var buf bytes.Buffer
			rlog.SetOutput(&buf)

			r := hermes.DefaultRouter()
			r.Use(middlewares.LoggingMiddleware)
			r.Get("/something", func(req hermes.Request, res hermes.Response) hermes.Result {
				return res.End()
			})
			r.Post("/something", func(req hermes.Request, res hermes.Response) hermes.Result {
				return res.End()
			})
			r.Put("/something/else", func(req hermes.Request, res hermes.Response) hermes.Result {
				return res.End()
			})

			handler := r.Handler()

			ctx1 := &fasthttp.RequestCtx{}
			ctx1.Request.Header.SetMethod("GET")
			ctx1.Request.URI().SetPath("/something")
			handler(ctx1)

			ctx2 := &fasthttp.RequestCtx{}
			ctx2.Request.Header.SetMethod("POST")
			ctx2.Request.URI().SetPath("/something")
			handler(ctx2)

			ctx3 := &fasthttp.RequestCtx{}
			ctx3.Request.Header.SetMethod("PUT")
			ctx3.Request.URI().SetPath("/something/else")
			handler(ctx3)

			logs := buf.String()
			Expect(logs).To(ContainSubstring("GET: /something                                         request_id=0"))
			Expect(logs).To(ContainSubstring("POST: /something                                         request_id=0"))
			Expect(logs).To(ContainSubstring("PUT: /something/else                                    request_id=0"))
		})
	})
})
