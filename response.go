package http

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var responsePool = &sync.Pool{
	New: func() interface{} {
		return &response{}
	},
}

type response struct {
	r *fasthttp.RequestCtx

	status int
}

func (res *response) reset() {
	res.r = nil
	res.status = 0
}

func acquireResponse(r *fasthttp.RequestCtx) *response {
	res := responsePool.Get().(*response)
	res.r = r
	return res
}

func releaseResponse(res *response) {
	res.reset()
	responsePool.Put(res)
}

func (res *response) Cookie(cookie *fasthttp.Cookie) Response {
	res.r.Response.Header.SetCookie(cookie)
	return res
}

func (res *response) Status(status int) Response {
	res.status = status
	return res
}

func (res *response) Header(name, value string) Response {
	res.r.Response.Header.Set(name, value)
	return res
}

func (res *response) Data(data interface{}) Result {
	return res.newResult().Data(data)
}

func (res *response) Error(err error) Result {
	return res.newResult().Error(err)
}

func (res *response) Redirect(uri string, code int) Result {
	return res.newResult().Redirect(uri, code)
}

func (res *response) File(filepath string) Result {
	return res.newResult().File(filepath)
}

func (res *response) FileDownload(filepath, filename string) Result {
	return res.newResult().FileDownload(filepath, filename)
}

func (res *response) End() Result {
	return res.newResult()
}

func (res *response) newResult() Result {
	return acquireResult(res.r, res.status)
}
