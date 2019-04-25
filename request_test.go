package http

import (
	"context"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/valyala/fasthttp"
)

var (
	describe  = g.Describe
	fdescribe = g.FDescribe
	it        = g.It
	fit       = g.FIt
)

func newRequest() *request {
	return &request{
		ctx: context.Background(),
		r:   &fasthttp.RequestCtx{},
	}
}

var _ = describe("Http", func() {
	describe("Request", func() {
		it("should accept the Content-type as JSON", func() {
			req := newRequest()
			req.Raw().Request.Header.SetContentType("application/json")
			req.Raw().Request.AppendBodyString(`{"foo":"bar"}`)
			Expect(req.IsJSON()).To(BeTrue())
		})

		it("should accept the Content-type with charset as JSON", func() {
			req := newRequest()
			req.Raw().Request.Header.SetContentType("application/json; charset=utf-8")
			req.Raw().Request.AppendBodyString(`{"foo":"bar"}`)
			Expect(req.IsJSON()).To(BeTrue())
		})

		it("should not accept the Content-type as JSON", func() {
			req := newRequest()
			Expect(req.IsJSON()).To(BeFalse())
		})

		it("should return the JSON structure in the body of the requisition", func() {
			req := newRequest()
			req.Raw().Request.Header.SetContentType("application/json; charset=utf-8")
			req.Raw().Request.AppendBodyString(`{"foo":"bar"}`)

			var data struct {
				Foo string `json:"foo"`
			}
			Expect(req.Data(&data)).To(BeNil())
			Expect(data.Foo).To(Equal("bar"))
		})
	})
})
