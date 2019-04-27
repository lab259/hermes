package http

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/lab259/errors"
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
		r: &fasthttp.RequestCtx{},
	}
}

var _ = describe("Http", func() {
	describe("Response", func() {
		it("should serialize and send a map as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(map[string]interface{}{
				"foo": "bar",
			})

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		it("should serialize and send a struct as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(simpleModel{
				Foo: "bar",
			})

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		it("should serialize and send a pointer as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data(&simpleModel{
				Foo: "bar",
			})

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"foo":"bar"}`))
		})

		it("should serialize and send a array/slice as JSON with the rightful 'Content-type' header", func() {
			res := newResponse()
			res.Data([]*simpleModel{
				&simpleModel{
					Foo: "bar",
				},
				&simpleModel{
					Foo: "baz",
				},
			})

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`[{"foo":"bar"},{"foo":"baz"}]`))
		})

		it("should fail serializing a JSON struct", func() {
			defer func() {
				err := recover()
				Expect(err).To(BeAssignableToTypeOf(&json.MarshalerError{}))
			}()

			res := newResponse()
			res.Data(&errornousJson{})
		})

		it("should write bytes to the context body", func() {
			res := newResponse()
			res.Data([]byte("this is a test"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		it("should write string to the context body", func() {
			res := newResponse()
			res.Data("this is a test")
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		it("should write a buff to the context body", func() {
			res := newResponse()
			res.Data(strings.NewReader("this is a test"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
		})

		it("should serialize internal server error", func() {
			res := newResponse()
			res.Error(errForced)

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(res.r.Response.StatusCode()).To(Equal(500))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"internal-server-error","message":"We encountered an internal error or misconfiguration and was unable to complete your request."}`))
		})

		it("should serialize wrapped error", func() {
			res := newResponse()
			res.Error(errors.Wrap(errForced, errors.Http(400), errors.Module("tests"), errors.Code("forced-error"), errors.Message("An error was forced.")))

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(res.r.Response.StatusCode()).To(Equal(400))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"forced-error","message":"An error was forced.","module":"tests"}`))
		})

		it("should override status with wrapped error", func() {
			res := newResponse()
			res.Status(403).Error(errors.Wrap(errForced, errors.Http(400), errors.Module("tests"), errors.Code("forced-error"), errors.Message("An error was forced.")))

			Expect(string(res.r.Response.Header.ContentType())).To(Equal("application/json; charset=utf-8"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(res.r.Response.StatusCode()).To(Equal(403))
			Expect(strings.TrimSpace(tmp.String())).To(Equal(`{"code":"forced-error","message":"An error was forced.","module":"tests"}`))
		})

		it("should redirect", func() {
			res := newResponse()
			res.Redirect("http://localhost:5000/redirect-test", 302)

			Expect(res.r.Response.StatusCode()).To(Equal(302))
			Expect(string(res.r.Response.Header.Peek("Location"))).To(Equal("http://localhost:5000/redirect-test"))
		})

		it("should set cookies", func() {
			res := newResponse()

			cookie := fasthttp.AcquireCookie()
			cookie.SetKey("session")
			cookie.SetValue("6194438f-f2f5-48b5-867b-1767b0f7d408")
			defer fasthttp.ReleaseCookie(cookie)

			res.Cookie(cookie).Data([]byte("this is a test"))

			Expect(res.r.Response.StatusCode()).To(Equal(200))
			cookieSetted := false
			res.r.Response.Header.VisitAllCookie(func(key []byte, value []byte) {
				if string(key) == "session" {
					cookieSetted = true
					Expect(string(value)).To(Equal("session=6194438f-f2f5-48b5-867b-1767b0f7d408"))
				}
			})
			Expect(cookieSetted).To(BeTrue())
		})

		it("should set header", func() {
			res := newResponse()
			res.Header("Content-Type", "awesome/test").Data([]byte("this is a test"))
			tmp := bytes.NewBufferString("")
			res.r.Response.BodyWriteTo(tmp)
			Expect(tmp.String()).To(Equal(`this is a test`))
			Expect(string(res.r.Response.Header.Peek("Content-Type"))).To(Equal("awesome/test"))
		})

		it("should send file", func() {
			res := newResponse()
			res.File("examples/files/sample.pdf")

			Expect(res.r.Response.StatusCode()).To(Equal(200))
			Expect(string(res.r.Response.Header.Peek("Content-Type"))).To(Equal("application/pdf"))
		})

		it("should send file (download)", func() {
			res := newResponse()
			res.FileDownload("examples/files/sample.pdf", "expected.pdf")

			Expect(res.r.Response.StatusCode()).To(Equal(200))
			Expect(string(res.r.Response.Header.Peek("Content-Type"))).To(Equal("application/pdf"))
			Expect(string(res.r.Response.Header.Peek("Content-Disposition"))).To(Equal("attachment; filename=expected.pdf"))
		})
	})
})
