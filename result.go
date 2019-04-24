package http

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/valyala/fasthttp"
)

type result struct {
	r *fasthttp.RequestCtx

	status      int
	hasSentData bool
}

func newResult(r *fasthttp.RequestCtx, status int) *result {
	return &result{r: r, status: status}
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
			r.r.SetContentTypeBytes(applicationJSON)
			e := json.NewEncoder(r.r.Response.BodyWriter())
			e.Encode(data)
		default:
			r.r.SetContentTypeBytes(textPlain)
			r.r.Response.AppendBodyString(fmt.Sprintf("%v", data))
		}
	}

	return r
}

func (r *result) Error(err error) Result {
	if err == nil {
		return r
	}

	r.defaultStatus(fasthttp.StatusInternalServerError)

	// TODO: integrate with lab259/errors
	return r.Data(map[string]interface{}{
		"message": err.Error(),
	})
}

func (r *result) setStatus() {
	if r.status == 0 {
		r.status = fasthttp.StatusOK
	}

	r.r.SetStatusCode(r.status)
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
