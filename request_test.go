package hermes

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

type contextKey string

var contextKeyTestID contextKey = "test:context:id"

var _ = describe("Hermes", func() {
	describe("Request", func() {
		it("should have a context", func() {
			req := newRequest()
			req.WithContext(context.WithValue(req.Context(), contextKeyTestID, "123"))
			Expect(req.Context().Value(contextKeyTestID)).To(Equal("123"))
		})

		it("should want JSON", func() {
			req := newRequest()
			req.Raw().Request.Header.Set("Accept", "application/json, text/html")
			Expect(req.WantsJSON()).To(BeTrue())
		})

		it("should not want JSON", func() {
			req := newRequest()
			req.Raw().Request.Header.Set("Accept", "text/html, application/xhtml+xml, application/xml")
			Expect(req.WantsJSON()).To(BeFalse())
		})

		it("should not want JSON if not first", func() {
			req := newRequest()
			req.Raw().Request.Header.Set("Accept", "text/html,application/json")
			Expect(req.WantsJSON()).To(BeFalse())
		})

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

		it("should get header", func() {
			req := newRequest()
			req.Raw().Request.Header.Set("X-Awesome", "Fully")
			Expect(string(req.Header("X-Awesome"))).To(Equal("Fully"))
		})

		it("should get query value", func() {
			req := newRequest()
			req.Raw().QueryArgs().Set("is_test", "true")

			Expect(string(req.Query("is_test"))).To(Equal("true"))
		})

		it("should get multiple query values", func() {
			req := newRequest()
			req.Raw().QueryArgs().Add("id", "1")
			req.Raw().QueryArgs().Add("id", "2")
			req.Raw().QueryArgs().Add("id", "3")

			values := req.QueryMulti("id")
			found := make([]string, len(values))
			for i, v := range values {
				found[i] = string(v)
			}
			Expect(found).To(ConsistOf("1", "2", "3"))
		})

		it("should get post value", func() {
			req := newRequest()
			req.Raw().PostArgs().Set("is_test", "true")

			Expect(string(req.Post("is_test"))).To(Equal("true"))
		})

		it("should get multiple post values", func() {
			req := newRequest()
			req.Raw().PostArgs().Add("id", "1")
			req.Raw().PostArgs().Add("id", "2")
			req.Raw().PostArgs().Add("id", "3")

			values := req.PostMulti("id")
			found := make([]string, len(values))
			for i, v := range values {
				found[i] = string(v)
			}
			Expect(found).To(ConsistOf("1", "2", "3"))
		})

		it("should get cookie", func() {
			req := newRequest()
			req.Raw().Request.Header.SetCookie("session", "6194438f-f2f5-48b5-867b-1767b0f7d408")
			Expect(string(req.Cookie("session"))).To(Equal("6194438f-f2f5-48b5-867b-1767b0f7d408"))
		})

		it("should get host", func() {
			req := newRequest()
			req.Raw().Request.SetHost("www.gijoe.io")
			Expect(string(req.Host())).To(Equal("www.gijoe.io"))
		})

		it("should get URI", func() {
			req := newRequest()
			req.Raw().Request.SetRequestURI("http://localhost:5000/v1/api")

			uri := req.URI()
			Expect(string(uri.Host())).To(Equal("localhost:5000"))
			Expect(string(uri.Path())).To(Equal("/v1/api"))
			Expect(string(uri.Scheme())).To(Equal("http"))
		})
	})
})
