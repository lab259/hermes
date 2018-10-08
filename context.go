package http

import (
	"bytes"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

var (
	StrApplicationJson = []byte("application/json")
	StrXForwardedFor   = []byte("X-Forwarded-For")
)

// Context is the wrapper structure on top of the fasthttp.
type Context struct {
	Ctx      *fasthttp.RequestCtx
	Response *fasthttp.Response
	Request  *fasthttp.Request
}

// NewContext returns a new `Context` from a `fasthttp.RequestCtx`.
func NewContext(ctx *fasthttp.RequestCtx) *Context {
	return &Context{Ctx: ctx, Response: &ctx.Response, Request: &ctx.Request}
}

// Write writes `buff` to the body response. It is an alias for the
// `fasthttp.RequestCtx.Write`.
func (ctx *Context) Write(buff []byte) (int, error) {
	return ctx.Ctx.Write(buff)
}

// UserValue is an alias for the `fasthttp.RequestCtx.UserValue` method.
func (ctx *Context) UserValue(name string) interface{} {
	return ctx.Ctx.UserValue(name)
}

// BodyJson tries to Unmarshal the body content of the request as a JSON file.
// The destiny of the unmarshaled data is placed into the passed `dst` pointer.
// If the unmarshaling process fails `BodyJson` will return an error.
func (ctx *Context) BodyJson(dst interface{}) error {
	return json.Unmarshal(ctx.Request.Body(), dst)
}

// IsJson checks if the Content-type of the request is `application/json`.
func (ctx *Context) IsJson() bool {
	ct := ctx.Request.Header.ContentType()
	laj := len(StrApplicationJson)
	return bytes.Equal(ct, StrApplicationJson) ||
		((laj < len(ct)) && bytes.Equal(ct[:laj], StrApplicationJson) && (ct[laj] == ';'))
}

// SendJson marshals the given `obj` and then uses the `SendJsonBytes` to
// deliver the information to the response.
func (ctx *Context) SendJson(obj interface{}) error {
	if jsonBytes, err := json.Marshal(obj); err == nil {
		return ctx.SendJsonBytes(jsonBytes)
	} else {
		return err
	}
}

// SendJsonBytes prepares the repsonse header to send the JSON Data. Afterwards,
// it appends the received `data` to the body.
func (ctx *Context) SendJsonBytes(data []byte) error {
	ctx.Response.Header.SetBytesV("Content-type", StrApplicationJson)
	ctx.Response.AppendBody(data)
	return nil
}
