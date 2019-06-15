package hermes

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/lab259/errors"
	"github.com/valyala/bytebufferpool"

	"github.com/valyala/fasthttp"
)

type result struct {
	r *fasthttp.RequestCtx

	status      int
	hasSentData bool
}

func (r *result) Data(data interface{}) Result {
	if r.hasSentData {
		return r
	}

	if err, ok := data.(error); ok {
		return r.Error(err)
	}

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
			if err := e.Encode(data); err != nil {
				return r.Error(err)
			}
		default:
			r.setContentType(defaultContentType)
			r.r.Response.AppendBodyString(fmt.Sprintf("%v", data))
		}
	}

	r.setStatus()
	r.hasSentData = true
	return r
}

func (r *result) Error(err error) Result {
	if err == nil {
		return r
	}

	errResponse := acquireErrorResponse(StatusInternalServerError)
	if !errors.AggregateToResponse(err, errResponse) {
		errResponse.SetParam("code", InternalServerErrorCode)
		errResponse.SetParam("message", InternalServerErrorMessage)
	}

	r.defaultStatus(errResponse.Status)
	r.Data(errResponse.Data)
	releaseErrorResponse(errResponse)
	return r
}

func (r *result) setStatus() {
	if r.status == 0 {
		r.status = StatusOK
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
	buff.SetString("attachment; filename=")
	buff.WriteString(filename)
	r.r.Response.Header.Set("Content-Disposition", buff.String())
	bytebufferpool.Put(buff)
	return r
}

func (r *result) End() {
	r.r = nil
	r.status = 0
	r.hasSentData = false
}
