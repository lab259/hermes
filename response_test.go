package hermes

import (
	"bytes"
	"strings"

	"github.com/lab259/errors/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/valyala/fasthttp"
)

var errForced = errors.New("forced error")

type simpleModel struct {
	Foo string `json:"foo"`
}

type errornousJson struct {
}

func (*errornousJson) MarshalJSON() ([]byte, error) {
	return nil, errForced
}

func newResponse() *response {
	return &response{
		result: result{r: &fasthttp.RequestCtx{}},
	}
}

var _ = Describe("Hermes", func() {
	Describe("Response", func() {
		It("should serialize and send a map as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(map[string]interface{}{
				"foo": "bar",
			})

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		It("should serialize and send a struct as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(simpleModel{
				Foo: "bar",
			})

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		It("should serialize and send a pointer as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(&simpleModel{
				Foo: "bar",
			})

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		It("should serialize and send a array/slice as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data([]*simpleModel{
				&simpleModel{
					Foo: "bar",
				},
				&simpleModel{
					Foo: "baz",
				},
			})

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`[{"foo":"bar"},{"foo":"baz"}]`))
		})

		It("should fail serializing a JSON struct", func() {
			res := newResponse()
			res.Data(&errornousJson{})
			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(res.result.r.Response.StatusCode()).To(Equal(500))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"internal-server-error","message":"We encountered an internal error or misconfiguration and was unable to complete your request."}`))
		})

		It("should write bytes to the context body", func() {
			res := newResponse()
			res.Data([]byte("this is a test"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		It("should write string to the context body", func() {
			res := newResponse()
			res.Data("this is a test")
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		It("should write a buff to the context body", func() {
			res := newResponse()
			res.Data(strings.NewReader("this is a test"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		It("should serialize internal server error", func() {
			res := newResponse()
			res.Error(errForced)

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(res.result.r.Response.StatusCode()).To(Equal(500))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"internal-server-error","message":"We encountered an internal error or misconfiguration and was unable to complete your request."}`))
		})

		It("should serialize wrapped error", func() {
			res := newResponse()
			res.Error(errors.Wrap(errForced, errors.Http(400), errors.Module("tests"), errors.Code("forced-error"), errors.Message("An error was forced.")))

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(res.result.r.Response.StatusCode()).To(Equal(400))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"forced-error","message":"An error was forced.","module":"tests"}`))
		})

		It("should override status with wrapped error", func() {
			res := newResponse()
			res.Status(403).Error(errors.Wrap(errForced, errors.Http(400), errors.Module("tests"), errors.Code("forced-error"), errors.Message("An error was forced.")))

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(res.result.r.Response.StatusCode()).To(Equal(403))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"forced-error","message":"An error was forced.","module":"tests"}`))
		})

		It("should wrap error with options", func() {
			res := newResponse()
			res.Status(403).Error(
				errForced,
				errors.Http(400),
				errors.Module("tests"),
				errors.Code("forced-error"),
				errors.Message("An error was forced."),
			)

			Expect(string(res.result.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(res.result.r.Response.StatusCode()).To(Equal(403))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"forced-error","message":"An error was forced.","module":"tests"}`))
		})

		It("should redirect", func() {
			res := newResponse()
			res.Redirect("http://localhost:5000/redirect-test", 302)

			Expect(res.result.r.Response.StatusCode()).To(Equal(302))
			Expect(string(res.result.r.Response.Header.Peek("Location"))).To(Equal("http://localhost:5000/redirect-test"))
		})

		It("should set cookies", func() {
			res := newResponse()

			cookie := fasthttp.AcquireCookie()
			cookie.SetKey("session")
			cookie.SetValue("6194438f-f2f5-48b5-867b-1767b0f7d408")
			defer fasthttp.ReleaseCookie(cookie)

			res.Cookie(cookie).Data([]byte("this is a test"))

			Expect(res.result.r.Response.StatusCode()).To(Equal(200))
			cookieSetted := false
			res.result.r.Response.Header.VisitAllCookie(func(key []byte, value []byte) {
				if string(key) == "session" {
					cookieSetted = true
					Expect(string(value)).To(Equal("session=6194438f-f2f5-48b5-867b-1767b0f7d408"))
				}
			})
			Expect(cookieSetted).To(BeTrue())
		})

		It("should set header", func() {
			res := newResponse()
			res.Header("Content-Type", "awesome/test").Data([]byte("this is a test"))
			tmp := bytes.NewBufferString("")
			res.result.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
			Expect(string(res.result.r.Response.Header.Peek("Content-Type"))).To(Equal("awesome/test"))
		})

		It("should send file", func() {
			res := newResponse()
			res.File("examples/files/sample.pdf")

			Expect(res.result.r.Response.StatusCode()).To(Equal(200))
			Expect(string(res.result.r.Response.Header.Peek("Content-Type"))).To(Equal("application/pdf"))
		})

		It("should send file (download)", func() {
			res := newResponse()
			res.FileDownload("examples/files/sample.pdf", "expected.pdf")

			Expect(res.result.r.Response.StatusCode()).To(Equal(200))
			Expect(string(res.result.r.Response.Header.Peek("Content-Type"))).To(Equal("application/pdf"))
			Expect(string(res.result.r.Response.Header.Peek("Content-Disposition"))).To(Equal("attachment; filename=expected.pdf"))
		})
	})
})
