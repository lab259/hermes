package hermes

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/valyala/fasthttp"
)

var requestPool = &sync.Pool{
	New: func() interface{} {
		return &BaseRequest{
			validParams: make([]string, 0, 10),
			params:      make([][]byte, 0, 10),
		}
	},
}

type BaseRequest struct {
	ctx         context.Context
	r           *fasthttp.RequestCtx
	validParams []string
	params      [][]byte
}

func AcquireRequest(ctx context.Context, r *fasthttp.RequestCtx) *BaseRequest {
	req := requestPool.Get().(*BaseRequest)
	req.r = r
	req.ctx = ctx
	return req
}

func ReleaseRequest(req *BaseRequest) {
	req.reset()
	requestPool.Put(req)
}

func (req *BaseRequest) reset() {
	req.r = nil
	req.ctx = nil
	req.validParams = req.validParams[:0]
	req.params = req.params[:0]
}

func (req *BaseRequest) Raw() *fasthttp.RequestCtx {
	return req.r
}

func (req *BaseRequest) Path() []byte {
	return req.r.Path()
}

func (req *BaseRequest) Method() []byte {
	return req.r.Method()
}

func (req *BaseRequest) URI() *fasthttp.URI {
	return req.r.URI()
}

func (req *BaseRequest) Header(name string) []byte {
	return req.r.Request.Header.Peek(name)
}

func (req *BaseRequest) Host() []byte {
	return req.r.Host()
}

func (req *BaseRequest) Param(name string) string {
	// req.params is not safe, since its reused over requests
	// but validParams is, so we check if name is one of the
	// valid params, before actually return the value
	for i, p := range req.validParams {
		if p == name {
			return string(req.params[i])
		}
	}
	return ""
}

func (req *BaseRequest) Query(name string) []byte {
	return req.r.QueryArgs().Peek(name)
}

func (req *BaseRequest) QueryMulti(name string) [][]byte {
	return req.r.QueryArgs().PeekMulti(name)
}

func (req *BaseRequest) Data(dst interface{}) error {
	return json.Unmarshal(req.r.PostBody(), dst)
}

func (req *BaseRequest) Post(name string) []byte {
	return req.r.PostArgs().Peek(name)
}

func (req *BaseRequest) PostMulti(name string) [][]byte {
	return req.r.PostArgs().PeekMulti(name)
}

func (req *BaseRequest) Cookie(name string) []byte {
	return req.r.Request.Header.Cookie(name)
}

func (req *BaseRequest) Context() context.Context {
	return req.ctx
}

func (req *BaseRequest) WithContext(ctx context.Context) Request {
	req.ctx = ctx
	return req
}

func (req *BaseRequest) IsJSON() bool {
	ct := req.r.Request.Header.ContentType()
	laj := len(applicationJSON)
	return bytes.Equal(ct, applicationJSON) ||
		((laj < len(ct)) && bytes.Equal(ct[:laj], applicationJSON) && (ct[laj] == ';'))
}

func (req *BaseRequest) WantsJSON() bool {
	accept := req.r.Request.Header.Peek("Accept")
	laj := len(applicationJSON)
	return bytes.Equal(accept, applicationJSON) ||
		((laj < len(accept)) && bytes.Equal(accept[:laj], applicationJSON))
}
