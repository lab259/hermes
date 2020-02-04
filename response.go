package hermes

import (
	"sync"

	"github.com/lab259/errors/v2"
	"github.com/valyala/fasthttp"
)

var responsePool = &sync.Pool{
	New: func() interface{} {
		return &BaseResponse{
			result: result{},
		}
	},
}

type BaseResponse struct {
	result result
}

func (res *BaseResponse) reset() {
	// result is resetted on .End()
}

func AcquireResponse(r *fasthttp.RequestCtx) *BaseResponse {
	res := responsePool.Get().(*BaseResponse)
	res.result.r = r
	return res
}

func ReleaseResponse(res *BaseResponse) {
	res.reset()
	responsePool.Put(res)
}

func (res *BaseResponse) Cookie(cookie *fasthttp.Cookie) Response {
	res.result.r.Response.Header.SetCookie(cookie)
	return res
}

func (res *BaseResponse) Status(status int) Response {
	res.result.status = status
	return res
}

func (res *BaseResponse) Header(name, value string) Response {
	res.result.r.Response.Header.Set(name, value)
	return res
}

func (res *BaseResponse) Data(data interface{}) Result {
	return res.result.Data(data)
}

func (res *BaseResponse) Error(err error, options ...interface{}) Result {
	return res.result.Error(errors.Wrap(err, options...))
}

func (res *BaseResponse) Redirect(uri string, code int) Result {
	return res.result.Redirect(uri, code)
}

func (res *BaseResponse) File(filepath string) Result {
	return res.result.File(filepath)
}

func (res *BaseResponse) FileDownload(filepath, filename string) Result {
	return res.result.FileDownload(filepath, filename)
}

func (res *BaseResponse) End() Result {
	return &res.result
}
