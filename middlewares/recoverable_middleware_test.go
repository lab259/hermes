package middlewares_test

import (
	"bytes"
	"strings"

	"github.com/lab259/http"
	"github.com/lab259/http/middlewares"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("Middlewares", func() {
	Describe("Recoverable Middleware", func() {
		It("should recover from panic", func() {
			r := http.DefaultRouter()
			r.Use(middlewares.RecoverableMiddleware)
			r.Get("/should-panic", func(req http.Request, res http.Response) http.Result {
				panic("The SDD alarm is down, bypass the redundant firewall so we can calculate the IB matrix!")
			})

			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod("GET")
			ctx.Request.URI().SetPath("/should-panic")

			r.Handler()(ctx)

			tmp := bytes.NewBufferString("")
			ctx.Response.BodyWriteTo(tmp)
			Expect(ctx.Response.StatusCode()).To(Equal(http.StatusInternalServerError))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"internal-server-error","message":"We encountered an internal error or misconfiguration and was unable to complete your request."}`))
		})
	})
})
