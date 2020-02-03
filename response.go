package hermes

import (
	"sync"

	"github.com/lab259/errors/v2"
	"github.com/valyala/fasthttp"
)

var responsePool = &sync.Pool{
	New: func() interface{} {
		return &response{
			result: result{},
		}
	},
}

type response struct {
	result result
}

func (res *response) reset() {
	// result is resetted on .End()
}

func AcquireResponse(r *fasthttp.RequestCtx) *response {
	res := responsePool.Get().(*response)
	res.result.r = r
	return res
}

func ReleaseResponse(res *response) {
	res.reset()
	responsePool.Put(res)
}

func (res *response) Cookie(cookie *fasthttp.Cookie) Response {
	res.result.r.Response.Header.SetCookie(cookie)
	return res
}

func (res *response) Status(status int) Response {
	res.result.status = status
	return res
}

func (res *response) Header(name, value string) Response {
	res.result.r.Response.Header.Set(name, value)
	return res
}

func (res *response) Data(data interface{}) Result {
	return res.result.Data(data)
}

func (res *response) Error(err error, options ...interface{}) Result {
	return res.result.Error(errors.Wrap(err, options...))
}

func (res *response) Redirect(uri string, code int) Result {
	return res.result.Redirect(uri, code)
}

func (res *response) File(filepath string) Result {
	return res.result.File(filepath)
}

func (res *response) FileDownload(filepath, filename string) Result {
	return res.result.FileDownload(filepath, filename)
}

func (res *response) End() Result {
	return &res.result
}
