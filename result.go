package http

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/lab259/errors"
	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fasthttp"
)

type result struct {
	r *fasthttp.RequestCtx

	status      int
	hasSentData bool
}

var resultPool = &sync.Pool{
	New: func() interface{} {
		return &result{}
	},
}

func acquireResult(r *fasthttp.RequestCtx, status int) *result {
	result := resultPool.Get().(*result)
	result.r = r
	result.status = status
	return result
}

func releaseResult(r *result) {
	resultPool.Put(r)
}

func (r *result) Data(data interface{}) Result {
	if r.hasSentData {
		return r
	}

	if err, ok := data.(error); ok {
		return r.Error(err)
	}

	r.hasSentData = true

	r.setStatus()
	if v, ok := data.([]byte); ok {
		r.r.Response.AppendBody(v)
	} else if v, ok := data.(io.Reader); ok {
		io.Copy(r.r.Response.BodyWriter(), v)
	} else {
		dataType := reflect.TypeOf(data)
		if dataType.Kind() == reflect.Ptr {
			dataType = dataType.Elem()
		}

		switch dataType.Kind() {
		case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
			r.setContentType(defaultJSONContentType)
			e := json.NewEncoder(r.r.Response.BodyWriter())
			err := e.Encode(data)
			if err != nil {
				panic(err)
			}
		default:
			r.setContentType(defaultContentType)
			r.r.Response.AppendBodyString(fmt.Sprintf("%v", data))
		}
	}

	return r
}

func (r *result) Error(err error) Result {
	if err == nil {
		return r
	}

	errResponse := acquireErrorResponse(fasthttp.StatusInternalServerError)
	defer releaseErrorResponse(errResponse)

	if !errors.AggregateToResponse(err, errResponse) {
		errResponse.SetParam("code", InternalServerErrorCode)
		errResponse.SetParam("message", InternalServerErrorMessage)
	}

	r.defaultStatus(errResponse.Status)
	return r.Data(errResponse.Data)
}

func (r *result) setStatus() {
	if r.status == 0 {
		r.status = fasthttp.StatusOK
	}

	r.r.SetStatusCode(r.status)
}

func (r *result) setContentType(v []byte) {
	r.r.SetContentTypeBytes(v)
}

func (r *result) defaultStatus(code int) {
	if r.status == 0 {
		r.status = code
	}
}

func (r *result) Redirect(uri string, code int) Result {
	r.r.Redirect(uri, code)
	return r
}

func (r *result) File(filepath string) Result {
	r.r.SendFile(filepath)
	return r
}

func (r *result) FileDownload(filepath, filename string) Result {
	r.r.SendFile(filepath)
	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	buff.SetString("attachment; filename=")
	buff.WriteString(filename)

	r.r.Response.Header.Set("Content-Disposition", buff.String())
	return r
}

func (r *result) End() {
	// reset and release
	r.r = nil
	r.status = 0
	r.hasSentData = false
	releaseResult(r)
}
