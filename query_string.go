package http

import (
	"strconv"
)

type QueryString interface {
	String(string, ...string) string
	Int(string, ...int) int
	Int64(string, ...int64) int64
	Float(string, ...float64) float64
	Bool(string, ...bool) bool
}

type queryString struct {
	req Request
}

func ParseQuery(req Request) QueryString {
	return &queryString{
		req: req,
	}
}

func (qs *queryString) String(s string, d ...string) string {
	v := qs.req.Query(s)
	if len(v) == 0 && len(d) > 0 {
		return d[0]
	}
	return string(v)
}

func (qs *queryString) Int(s string, d ...int) int {
	var defaultValue int

	v := qs.req.Query(s)
	if i, err := strconv.Atoi(string(v)); err == nil {
		return i
	}

	if len(d) > 0 {
		return d[0]
	}

	return defaultValue
}

func (qs *queryString) Int64(s string, d ...int64) int64 {
	var defaultValue int64

	v := qs.req.Query(s)
	if i, err := strconv.ParseInt(string(v), 10, 64); err == nil {
		return i
	}

	if len(d) > 0 {
		return d[0]
	}

	return defaultValue
}

func (qs *queryString) Float(s string, d ...float64) float64 {
	var defaultValue float64

	v := qs.req.Query(s)
	if f, err := strconv.ParseFloat(string(v), 64); err == nil {
		return f
	}

	if len(d) > 0 {
		return d[0]
	}

	return defaultValue
}

func (qs *queryString) Bool(s string, d ...bool) bool {
	v := qs.req.Query(s)
	if b, err := strconv.ParseBool(string(v)); err == nil {
		return b
	}

	if len(d) > 0 {
		return d[0]
	}

	return false
}
