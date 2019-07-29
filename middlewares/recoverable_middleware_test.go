package middlewares_test

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/lab259/errors/v2"
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("Middlewares", func() {
	Describe("Recoverable Middleware", func() {
		It("should recover from panic", func() {
			r := hermes.DefaultRouter()
			r.Use(middlewares.RecoverableMiddleware)
			r.Get("/should-panic", func(req hermes.Request, res hermes.Response) hermes.Result {
				panic("The SDD alarm is down, bypass the redundant firewall so we can calculate the IB matrix!")
			})

			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod("GET")
			ctx.Request.URI().SetPath("/should-panic")

			r.Handler()(ctx)

			tmp := bytes.NewBufferString("")
			ctx.Response.BodyWriteTo(tmp)
			Expect(ctx.Response.StatusCode()).To(Equal(hermes.StatusInternalServerError))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"internal-server-error","message":"We encountered an internal error or misconfiguration and was unable to complete your request."}`))
		})

		It("should recover from panic (formatted error)", func() {
			r := hermes.DefaultRouter()
			r.Use(middlewares.RecoverableMiddleware)
			r.Get("/should-panic", func(req hermes.Request, res hermes.Response) hermes.Result {
				panic(errors.Wrap(errors.New("forced-error"), "Something really bad happened.", hermes.StatusBadRequest, errors.Code("forced-error"), errors.Module("testing")))
			})

			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod("GET")
			ctx.Request.URI().SetPath("/should-panic")

			r.Handler()(ctx)

			tmp := bytes.NewBufferString("")
			ctx.Response.BodyWriteTo(tmp)
			Expect(ctx.Response.StatusCode()).To(Equal(hermes.StatusBadRequest))

			var errData map[string]interface{}
			Expect(json.Unmarshal(tmp.Bytes(), &errData))
			Expect(errData["code"]).To(Equal("forced-error"))
			Expect(errData["message"]).To(Equal("Something really bad happened."))
			Expect(errData["module"]).To(Equal("testing"))
		})
	})
})
