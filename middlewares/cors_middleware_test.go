package middlewares_test

import (
	"regexp"
	"strings"

	"github.com/lab259/cors"
	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("Middlewares", func() {
	Describe("CORS Middleware", func() {
		var allHeaders = []string{
			"Vary",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Credentials",
			"Access-Control-Max-Age",
			"Access-Control-Expose-Headers",
		}

		cases := []struct {
			name       string
			options    cors.Options
			method     string
			reqHeaders map[string]string
			resHeaders map[string]string
		}{
			{
				"Default",
				cors.Options{
					// Intentionally left blank.
				},
				"GET",
				map[string]string{},
				map[string]string{
					"Vary": "Origin",
				},
			},
			{
				"NoConfig",
				cors.Options{
					// Intentionally left blank.
				},
				"GET",
				map[string]string{},
				map[string]string{
					"Vary": "Origin",
				},
			},
			{
				"MatchAllOrigin",
				cors.Options{
					AllowedOrigins: []string{"*"},
				},
				"GET",
				map[string]string{
					"Origin": "http://foobar.com",
				},
				map[string]string{
					"Vary":                        "Origin",
					"Access-Control-Allow-Origin": "*",
				},
			},
			{
				"MatchAllOriginWithCredentials",
				cors.Options{
					AllowedOrigins:   []string{"*"},
					AllowCredentials: true,
				},
				"GET",
				map[string]string{
					"Origin": "http://foobar.com",
				},
				map[string]string{
					"Vary":                             "Origin",
					"Access-Control-Allow-Origin":      "*",
					"Access-Control-Allow-Credentials": "true",
				},
			},
			{
				"AllowedOrigin",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
				},
				"GET",
				map[string]string{
					"Origin": "http://foobar.com",
				},
				map[string]string{
					"Vary":                        "Origin",
					"Access-Control-Allow-Origin": "http://foobar.com",
				},
			},
			{
				"WildcardOrigin",
				cors.Options{
					AllowedOrigins: []string{"http://*.bar.com"},
				},
				"GET",
				map[string]string{
					"Origin": "http://foo.bar.com",
				},
				map[string]string{
					"Vary":                        "Origin",
					"Access-Control-Allow-Origin": "http://foo.bar.com",
				},
			},
			{
				"DisallowedOrigin",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
				},
				"GET",
				map[string]string{
					"Origin": "http://barbaz.com",
				},
				map[string]string{
					"Vary": "Origin",
				},
			},
			{
				"DisallowedWildcardOrigin",
				cors.Options{
					AllowedOrigins: []string{"http://*.bar.com"},
				},
				"GET",
				map[string]string{
					"Origin": "http://foo.baz.com",
				},
				map[string]string{
					"Vary": "Origin",
				},
			},
			{
				"AllowedOriginFuncMatch",
				cors.Options{
					AllowOriginFunc: func(o string) bool {
						return regexp.MustCompile("^http://foo").MatchString(o)
					},
				},
				"GET",
				map[string]string{
					"Origin": "http://foobar.com",
				},
				map[string]string{
					"Vary":                        "Origin",
					"Access-Control-Allow-Origin": "http://foobar.com",
				},
			},
			{
				"AllowOriginRequestFuncMatch",
				cors.Options{
					AllowOriginRequestFunc: func(ctx *fasthttp.RequestCtx, o string) bool {
						return regexp.MustCompile("^http://foo").MatchString(o) && string(ctx.Request.Header.Peek("Authorization")) == "secret"
					},
				},
				"GET",
				map[string]string{
					"Origin":        "http://foobar.com",
					"Authorization": "secret",
				},
				map[string]string{
					"Vary":                        "Origin",
					"Access-Control-Allow-Origin": "http://foobar.com",
				},
			},
			{
				"AllowOriginRequestFuncNotMatch",
				cors.Options{
					AllowOriginRequestFunc: func(ctx *fasthttp.RequestCtx, o string) bool {
						return regexp.MustCompile("^http://foo").MatchString(o) && string(ctx.Request.Header.Peek("Authorization")) == "secret"
					},
				},
				"GET",
				map[string]string{
					"Origin":        "http://foobar.com",
					"Authorization": "not-secret",
				},
				map[string]string{
					"Vary": "Origin",
				},
			},
			{
				"MaxAge",
				cors.Options{
					AllowedOrigins: []string{"http://example.com/"},
					AllowedMethods: []string{"GET"},
					MaxAge:         10,
				},
				"OPTIONS",
				map[string]string{
					"Origin":                        "http://example.com/",
					"Access-Control-Request-Method": "GET",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://example.com/",
					"Access-Control-Allow-Methods": "GET",
					"Access-Control-Max-Age":       "10",
				},
			},
			{
				"AllowedMethod",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedMethods: []string{"PUT", "DELETE"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                        "http://foobar.com",
					"Access-Control-Request-Method": "PUT",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://foobar.com",
					"Access-Control-Allow-Methods": "PUT",
				},
			},
			{
				"DisallowedMethod",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedMethods: []string{"PUT", "DELETE"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                        "http://foobar.com",
					"Access-Control-Request-Method": "PATCH",
				},
				map[string]string{
					"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				},
			},
			{
				"AllowedHeaders",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedHeaders: []string{"X-Header-1", "x-header-2"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                         "http://foobar.com",
					"Access-Control-Request-Method":  "GET",
					"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://foobar.com",
					"Access-Control-Allow-Methods": "GET",
					"Access-Control-Allow-Headers": "X-Header-2, X-Header-1",
				},
			},
			{
				"DefaultAllowedHeaders",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedHeaders: []string{},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                         "http://foobar.com",
					"Access-Control-Request-Method":  "GET",
					"Access-Control-Request-Headers": "X-Requested-With",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://foobar.com",
					"Access-Control-Allow-Methods": "GET",
					"Access-Control-Allow-Headers": "X-Requested-With",
				},
			},
			{
				"AllowedWildcardHeader",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedHeaders: []string{"*"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                         "http://foobar.com",
					"Access-Control-Request-Method":  "GET",
					"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://foobar.com",
					"Access-Control-Allow-Methods": "GET",
					"Access-Control-Allow-Headers": "X-Header-2, X-Header-1",
				},
			},
			{
				"DisallowedHeader",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					AllowedHeaders: []string{"X-Header-1", "x-header-2"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                         "http://foobar.com",
					"Access-Control-Request-Method":  "GET",
					"Access-Control-Request-Headers": "X-Header-3, X-Header-1",
				},
				map[string]string{
					"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				},
			},
			{
				"OriginHeader",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
				},
				"OPTIONS",
				map[string]string{
					"Origin":                         "http://foobar.com",
					"Access-Control-Request-Method":  "GET",
					"Access-Control-Request-Headers": "origin",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "http://foobar.com",
					"Access-Control-Allow-Methods": "GET",
					"Access-Control-Allow-Headers": "Origin",
				},
			},
			{
				"ExposedHeader",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
					ExposedHeaders: []string{"X-Header-1", "x-header-2"},
				},
				"GET",
				map[string]string{
					"Origin": "http://foobar.com",
				},
				map[string]string{
					"Vary":                          "Origin",
					"Access-Control-Allow-Origin":   "http://foobar.com",
					"Access-Control-Expose-Headers": "X-Header-1, X-Header-2",
				},
			},
			{
				"AllowedCredentials",
				cors.Options{
					AllowedOrigins:   []string{"http://foobar.com"},
					AllowCredentials: true,
				},
				"OPTIONS",
				map[string]string{
					"Origin":                        "http://foobar.com",
					"Access-Control-Request-Method": "GET",
				},
				map[string]string{
					"Vary":                             "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":      "http://foobar.com",
					"Access-Control-Allow-Methods":     "GET",
					"Access-Control-Allow-Credentials": "true",
				},
			},
			{
				"OptionPassthrough",
				cors.Options{
					OptionsPassthrough: true,
				},
				"OPTIONS",
				map[string]string{
					"Origin":                        "http://foobar.com",
					"Access-Control-Request-Method": "GET",
				},
				map[string]string{
					"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET",
				},
			},
			{
				"NonPreflightOptions",
				cors.Options{
					AllowedOrigins: []string{"http://foobar.com"},
				},
				"OPTIONS",
				map[string]string{},
				map[string]string{},
			},
		}
		for i := range cases {
			tc := cases[i]
			It(tc.name, func() {
				r := hermes.DefaultRouter()

				if tc.name == "Default" {
					r.Use(middlewares.DefaultCorsMiddleware())
				} else {
					r.Use(middlewares.NewCorsMiddleware(tc.options))
				}

				r.Get("/something", func(req hermes.Request, res hermes.Response) hermes.Result {
					return res.End()
				})

				ctx := fasthttp.RequestCtx{}
				ctx.Request.Header.SetMethod(tc.method)
				ctx.Request.SetRequestURI("http://example.com/foo")
				for name, value := range tc.reqHeaders {
					ctx.Request.Header.Add(name, value)
				}

				r.Handler()(&ctx)

				headers := make(map[string][]string, ctx.Response.Header.Len())
				ctx.Response.Header.VisitAll(func(key []byte, value []byte) {
					if arr, found := headers[string(key)]; found {
						headers[string(key)] = append(arr, string(value))
					} else {
						headers[string(key)] = []string{string(value)}
					}
				})

				for _, name := range allHeaders {
					got := strings.Join(headers[name], ", ")
					want := tc.resHeaders[name]
					Expect(got).To(Equal(want))
				}
			})
		}
	})
})
