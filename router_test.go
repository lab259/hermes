package http

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"github.com/valyala/fasthttp"
	"testing"
)

func createRequestCtxFromPath(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.URI().SetPath(path)
	return ctx
}

var emptyHandler = func(ctx *Context) {
}

var _ = ginkgo.Describe("Router", func() {

	ginkgo.Describe("Split", func() {
		ginkgo.It("should split the path", func() {
			path := []byte("/path/with/four/parts")
			tokens := make([][]byte, 0)
			tokens = Split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		ginkgo.It("should split the path not starting with /", func() {
			path := []byte("path/with/four/parts")
			tokens := make([][]byte, 0)
			tokens = Split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		ginkgo.It("should split the path ending with /", func() {
			path := []byte("/path/with/four/parts/")
			tokens := make([][]byte, 0)
			tokens = Split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		ginkgo.It("should split an empty path", func() {
			path := []byte("/")
			tokens := make([][]byte, 0)
			tokens = Split(path, tokens)
			Expect(tokens).To(BeEmpty())
		})
	})

	ginkgo.Describe("Parse", func() {

		ginkgo.It("should parse a GET", func() {
			router := NewRouter()
			router.GET("", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].handler).NotTo(BeNil())
		})

		ginkgo.It("should parse a GET", func() {
			router := NewRouter()
			router.GET("/route", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("route"))
			Expect(router.children["GET"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a POST", func() {
			router := NewRouter()
			router.POST("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a PUT", func() {
			router := NewRouter()
			router.PUT("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PUT"))
			Expect(router.children["PUT"].children).To(HaveKey("route"))
			Expect(router.children["PUT"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a DELETE", func() {
			router := NewRouter()
			router.DELETE("/route", emptyHandler)

			Expect(router.children).To(HaveKey("DELETE"))
			Expect(router.children["DELETE"].children).To(HaveKey("route"))
			Expect(router.children["DELETE"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a HEAD", func() {
			router := NewRouter()
			router.HEAD("/route", emptyHandler)

			Expect(router.children).To(HaveKey("HEAD"))
			Expect(router.children["HEAD"].children).To(HaveKey("route"))
			Expect(router.children["HEAD"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a OPTIONS", func() {
			router := NewRouter()
			router.OPTIONS("/route", emptyHandler)

			Expect(router.children).To(HaveKey("OPTIONS"))
			Expect(router.children["OPTIONS"].children).To(HaveKey("route"))
			Expect(router.children["OPTIONS"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a PATCH", func() {
			router := NewRouter()
			router.PATCH("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PATCH"))
			Expect(router.children["PATCH"].children).To(HaveKey("route"))
			Expect(router.children["PATCH"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a POST", func() {
			router := NewRouter()
			router.POST("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		ginkgo.It("should parse a complete static route", func() {
			router := NewRouter()
			router.GET("/this/should/be/static", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("this"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children).To(HaveKey("should"))
			Expect(router.children["GET"].children["this"].children["should"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this"].children["should"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))
		})

		ginkgo.It("should parse multiple static routes related", func() {
			router := NewRouter()
			router.GET("/this/should/be/static", emptyHandler)
			router.GET("/this/should2/be/static", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("this"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].handler).To(BeNil())

			Expect(router.children["GET"].children["this"].children).To(HaveKey("should"))
			Expect(router.children["GET"].children["this"].children["should"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this"].children["should"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))

			Expect(router.children["GET"].children["this"].children).To(HaveKey("should2"))
			Expect(router.children["GET"].children["this"].children["should2"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this"].children["should2"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))
		})

		ginkgo.It("should parse a complete a route starting static and ending with a wildcard", func() {
			router := NewRouter()
			router.GET("/static/:wildcard", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].names).To(BeEmpty())
			Expect(router.children["GET"].children).To(HaveKey("static"))
			Expect(router.children["GET"].children["static"].children).To(BeEmpty())
			Expect(router.children["GET"].children["static"].handler).To(BeNil())
			Expect(router.children["GET"].children["static"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].children["static"].wildcard.handler).NotTo(BeNil())
			Expect(router.children["GET"].children["static"].wildcard.children).To(BeEmpty())
			Expect(router.children["GET"].children["static"].wildcard.names).To(Equal([]string{"wildcard"}))
		})

		ginkgo.It("should parse multiple static routes related and not", func() {
			router := NewRouter()
			router.GET("/this/should/be/static", emptyHandler)
			router.GET("/this/should2/be/static", emptyHandler)
			router.GET("/this2/should/be/static", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("this"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children).To(HaveKey("should"))
			Expect(router.children["GET"].children["this"].children["should"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this"].children["should"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))

			Expect(router.children["GET"].children["this"].children).To(HaveKey("should2"))
			Expect(router.children["GET"].children["this"].children["should2"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this"].children["should2"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))

			Expect(router.children["GET"].children).To(HaveKey("this2"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this2"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this2"].handler).To(BeNil())
			Expect(router.children["GET"].children["this2"].children).To(HaveKey("should"))
			Expect(router.children["GET"].children["this2"].children["should"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this2"].children["should"].handler).To(BeNil())
			Expect(router.children["GET"].children["this2"].children["should"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this2"].children["should"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this2"].children["should"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this2"].children["should"].children["be"].children).To(HaveKey("static"))
			Expect(fmt.Sprintf("%p", router.children["GET"].children["this2"].children["should"].children["be"].children["static"].handler)).To(Equal(fmt.Sprintf("%p", emptyHandler)))
		})

		ginkgo.It("should parse a complete route with wildcard", func() {
			router := NewRouter()
			router.GET("/:account/detail/another", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.children).To(HaveKey("detail"))
			Expect(router.children["GET"].wildcard.handler).To(BeNil())
			Expect(router.children["GET"].wildcard.wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].children).To(HaveKey("another"))
			Expect(router.children["GET"].wildcard.children["detail"].handler).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].children["another"].wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].children["another"].handler).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].children["another"].names).To(Equal([]string{"account"}))
		})

		ginkgo.It("should parse a complete route with a sequence of wildcards", func() {
			router := NewRouter()
			router.GET("/:account/:transaction/:invoice", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.handler).To(BeNil())
			Expect(router.children["GET"].wildcard.names).To(BeEmpty())
			Expect(router.children["GET"].wildcard.wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.wildcard.handler).To(BeNil())
			Expect(router.children["GET"].wildcard.wildcard.names).To(BeEmpty())
			Expect(router.children["GET"].wildcard.wildcard.wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.wildcard.wildcard.wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.wildcard.wildcard.handler).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.wildcard.wildcard.names).To(Equal([]string{"account", "transaction", "invoice"}))
		})

		ginkgo.It("should parse multiple routes starting with wildcards", func() {
			router := NewRouter()
			router.GET("/:account/detail", emptyHandler)
			router.GET("/:account/history", emptyHandler)
			router.GET("/:transaction/invoice", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.children).To(HaveKey("detail"))
			Expect(router.children["GET"].wildcard.handler).To(BeNil())
			Expect(router.children["GET"].wildcard.wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard.children["detail"].handler).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["detail"].names).To(Equal([]string{"account"}))
			Expect(router.children["GET"].wildcard.children["history"].names).To(Equal([]string{"account"}))
			Expect(router.children["GET"].wildcard.children["invoice"].names).To(Equal([]string{"transaction"}))
		})

		ginkgo.It("should parse multiple mixed routes", func() {
			router := NewRouter()
			router.GET("/accounts/:account/detail", emptyHandler)
			router.GET("/accounts/:account/history", emptyHandler)
			router.GET("/:transaction/invoice", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveLen(1))
			Expect(router.children["GET"].children).To(HaveKey("accounts"))
			Expect(router.children["GET"].children["accounts"].children).To(BeEmpty())
			Expect(router.children["GET"].children["accounts"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].children["accounts"].handler).To(BeNil())
			Expect(router.children["GET"].children["accounts"].wildcard.children).To(HaveLen(2))
			Expect(router.children["GET"].children["accounts"].wildcard.children).To(HaveKey("detail"))
			Expect(router.children["GET"].children["accounts"].wildcard.children["detail"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["accounts"].wildcard.children["detail"].children).To(BeEmpty())
			Expect(router.children["GET"].children["accounts"].wildcard.children["detail"].handler).NotTo(BeNil())
			Expect(router.children["GET"].children["accounts"].wildcard.children["detail"].names).To(Equal([]string{"account"}))
			Expect(router.children["GET"].children["accounts"].wildcard.children).To(HaveKey("history"))
			Expect(router.children["GET"].children["accounts"].wildcard.children["history"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["accounts"].wildcard.children["history"].children).To(BeEmpty())
			Expect(router.children["GET"].children["accounts"].wildcard.children["history"].handler).NotTo(BeNil())
			Expect(router.children["GET"].children["accounts"].wildcard.children["history"].names).To(Equal([]string{"account"}))
			Expect(router.children["GET"].wildcard).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.handler).To(BeNil())
			Expect(router.children["GET"].wildcard.children).To(HaveKey("invoice"))
			Expect(router.children["GET"].wildcard.children["invoice"].wildcard).To(BeNil())
			Expect(router.children["GET"].wildcard.children["invoice"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard.children["invoice"].handler).NotTo(BeNil())
			Expect(router.children["GET"].wildcard.children["invoice"].names).To(Equal([]string{"transaction"}))
		})

		ginkgo.It("should panic due to conflicting empty tokens", func() {
			router := NewRouter()

			Expect(func() {
				router.GET("//detail", emptyHandler)
			}).To(Panic())

			Expect(func() {
				router.GET("/account/detail//", emptyHandler)
			}).To(Panic())

			Expect(func() {
				router.GET("/account//detail", emptyHandler)
			}).To(Panic())
		})

		ginkgo.It("should not panic with empty token at the end", func() {
			router := NewRouter()

			Expect(func() {
				router.GET("/account/", emptyHandler)
			}).NotTo(Panic())

			Expect(func() {
				router.GET("/account/detail/", emptyHandler)
			}).NotTo(Panic())

			Expect(func() {
				router.GET("/account/detail/:id/", emptyHandler)
			}).NotTo(Panic())
		})

		ginkgo.It("should panic due to conflicting static routes", func() {
			router := NewRouter()
			router.GET("/account/detail", emptyHandler)
			Expect(func() {
				router.GET("/account/detail", emptyHandler)
			}).To(Panic())
		})

		ginkgo.It("should panic due to conflicting 'wildcarded' routes", func() {
			router := NewRouter()
			router.GET("/:account", emptyHandler)
			Expect(func() {
				router.GET("/:transaction", emptyHandler)
			}).To(Panic())
		})

		ginkgo.It("should panic due to conflicting mixing routes", func() {
			router := NewRouter()
			router.GET("/:account/detail", emptyHandler)
			router.GET("/:account/id", emptyHandler)
			Expect(func() {
				router.GET("/:transaction/id", emptyHandler)
			}).To(Panic())
		})

		ginkgo.It("should not match any ropute", func() {
			router := NewRouter()
			router.GET("/:account/detail", emptyHandler)
			router.GET("/:account/id", emptyHandler)
			ok, _, _ := router.children["GET"].Matches(nil, nil)
			Expect(ok).To(BeFalse())
		})
	})

	ginkgo.Describe("Group", func() {
		ginkgo.It("should parse a GET", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.GET("/route", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].children).To(HaveLen(1))
			Expect(router.children["GET"].children).To(HaveKey("group"))
			Expect(router.children["GET"].children["group"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["group"].handler).To(BeNil())
			Expect(router.children["GET"].children["group"].children).To(HaveLen(1))
			Expect(router.children["GET"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["GET"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["GET"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a POST", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.POST("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].wildcard).To(BeNil())
			Expect(router.children["POST"].children).To(HaveLen(1))
			Expect(router.children["POST"].children).To(HaveKey("group"))
			Expect(router.children["POST"].children["group"].wildcard).To(BeNil())
			Expect(router.children["POST"].children["group"].handler).To(BeNil())
			Expect(router.children["POST"].children["group"].children).To(HaveLen(1))
			Expect(router.children["POST"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["POST"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["POST"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["POST"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a PUT", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.PUT("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PUT"))
			Expect(router.children["PUT"].wildcard).To(BeNil())
			Expect(router.children["PUT"].children).To(HaveLen(1))
			Expect(router.children["PUT"].children).To(HaveKey("group"))
			Expect(router.children["PUT"].children["group"].wildcard).To(BeNil())
			Expect(router.children["PUT"].children["group"].handler).To(BeNil())
			Expect(router.children["PUT"].children["group"].children).To(HaveLen(1))
			Expect(router.children["PUT"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["PUT"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["PUT"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["PUT"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a DELETE", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.DELETE("/route", emptyHandler)

			Expect(router.children).To(HaveKey("DELETE"))
			Expect(router.children["DELETE"].wildcard).To(BeNil())
			Expect(router.children["DELETE"].children).To(HaveLen(1))
			Expect(router.children["DELETE"].children).To(HaveKey("group"))
			Expect(router.children["DELETE"].children["group"].wildcard).To(BeNil())
			Expect(router.children["DELETE"].children["group"].handler).To(BeNil())
			Expect(router.children["DELETE"].children["group"].children).To(HaveLen(1))
			Expect(router.children["DELETE"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["DELETE"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["DELETE"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["DELETE"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a HEAD", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.HEAD("/route", emptyHandler)

			Expect(router.children).To(HaveKey("HEAD"))
			Expect(router.children["HEAD"].wildcard).To(BeNil())
			Expect(router.children["HEAD"].children).To(HaveLen(1))
			Expect(router.children["HEAD"].children).To(HaveKey("group"))
			Expect(router.children["HEAD"].children["group"].wildcard).To(BeNil())
			Expect(router.children["HEAD"].children["group"].handler).To(BeNil())
			Expect(router.children["HEAD"].children["group"].children).To(HaveLen(1))
			Expect(router.children["HEAD"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["HEAD"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["HEAD"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["HEAD"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a OPTIONS", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.OPTIONS("/route", emptyHandler)

			Expect(router.children).To(HaveKey("OPTIONS"))
			Expect(router.children["OPTIONS"].wildcard).To(BeNil())
			Expect(router.children["OPTIONS"].children).To(HaveLen(1))
			Expect(router.children["OPTIONS"].children).To(HaveKey("group"))
			Expect(router.children["OPTIONS"].children["group"].wildcard).To(BeNil())
			Expect(router.children["OPTIONS"].children["group"].handler).To(BeNil())
			Expect(router.children["OPTIONS"].children["group"].children).To(HaveLen(1))
			Expect(router.children["OPTIONS"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["OPTIONS"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["OPTIONS"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["OPTIONS"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should parse a PATCH", func() {
			router := NewRouter()
			group := router.Group("/group")
			group.PATCH("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PATCH"))
			Expect(router.children["PATCH"].wildcard).To(BeNil())
			Expect(router.children["PATCH"].children).To(HaveLen(1))
			Expect(router.children["PATCH"].children).To(HaveKey("group"))
			Expect(router.children["PATCH"].children["group"].wildcard).To(BeNil())
			Expect(router.children["PATCH"].children["group"].handler).To(BeNil())
			Expect(router.children["PATCH"].children["group"].children).To(HaveLen(1))
			Expect(router.children["PATCH"].children["group"].children).To(HaveKey("route"))
			Expect(router.children["PATCH"].children["group"].children["route"].wildcard).To(BeNil())
			Expect(router.children["PATCH"].children["group"].children["route"].handler).NotTo(BeNil())
			Expect(router.children["PATCH"].children["group"].children["route"].children).To(BeEmpty())
		})

		ginkgo.It("should check the subgroup", func() {
			router := NewRouter()
			group := router.Group("/group")
			group2 := group.Group("/subgroup").(*routerGroup)

			Expect(group2).NotTo(BeNil())
			Expect(group2.router).To(Equal(group))
			Expect(group2.prefix).To(Equal("/subgroup"))
		})
	})

	ginkgo.Describe("Handle", func() {
		var router *Router

		ginkgo.BeforeEach(func() {
			router = NewRouter()
		})

		ginkgo.It("should resolve an empty route", func() {
			value := 1
			router.GET("", func(ctx *Context) {
				value = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/"))

			Expect(value).To(Equal(2))
		})

		ginkgo.It("should resolve an empty trailing route", func() {
			value := 1
			router.GET("/", func(ctx *Context) {
				value = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/"))

			Expect(value).To(Equal(2))
		})

		ginkgo.It("should resolve a static route", func() {
			value := 1
			router.GET("/static", func(ctx *Context) {
				value = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/static"))

			Expect(value).To(Equal(2))
		})

		ginkgo.It("should resolve a static route not starting with /", func() {
			value := 1
			router.GET("static", func(ctx *Context) {
				value = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/static"))

			Expect(value).To(Equal(2))
		})

		ginkgo.It("should resolve multiple static routes", func() {
			value1 := 1
			value2 := 1
			value3 := 1

			router.GET("/static", func(ctx *Context) {
				value1 = 2
			})

			router.GET("/static/second", func(ctx *Context) {
				value2 = 2
			})

			router.GET("/another", func(ctx *Context) {
				value3 = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/static"))
			router.Handler(createRequestCtxFromPath("GET", "/static/second"))
			router.Handler(createRequestCtxFromPath("GET", "/another"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		ginkgo.It("should resolve a wildcard route", func() {
			value := 1
			router.GET("/:wildcard", func(ctx *Context) {
				Expect(ctx.UserValue("wildcard")).To(Equal("value"))
				value = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/value"))

			Expect(value).To(Equal(2))
		})

		ginkgo.It("should resolve a multiple wildcard routes", func() {
			value1 := 1
			value2 := 1
			value3 := 1
			router.GET("/:account/transactions", func(ctx *Context) {
				Expect(ctx.UserValue("account")).To(Equal("value1"))
				value1 = 2
			})
			router.GET("/:account/profile", func(ctx *Context) {
				Expect(ctx.UserValue("account")).To(Equal("value2"))
				value2 = 2
			})
			router.GET("/:user/roles", func(ctx *Context) {
				Expect(ctx.UserValue("user")).To(Equal("value3"))
				value3 = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/value1/transactions"))
			router.Handler(createRequestCtxFromPath("GET", "/value2/profile"))
			router.Handler(createRequestCtxFromPath("GET", "/value3/roles"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		ginkgo.It("should resolve a multiple wildcard routes in sequence", func() {
			value1 := 1
			value2 := 1
			value3 := 1
			router.GET("/:account/:subscription/cancel", func(ctx *Context) {
				Expect(ctx.UserValue("account")).To(Equal("account1"))
				Expect(ctx.UserValue("subscription")).To(Equal("subscription1"))
				value1 = 2
			})
			router.GET("/:account/:subscription/history", func(ctx *Context) {
				Expect(ctx.UserValue("account")).To(Equal("account2"))
				Expect(ctx.UserValue("subscription")).To(Equal("subscription2"))
				value2 = 2
			})
			router.GET("/:account/:subscription", func(ctx *Context) {
				Expect(ctx.UserValue("account")).To(Equal("account3"))
				Expect(ctx.UserValue("subscription")).To(Equal("subscription3"))
				value3 = 2
			})

			router.Handler(createRequestCtxFromPath("GET", "/account1/subscription1/cancel"))
			router.Handler(createRequestCtxFromPath("GET", "/account2/subscription2/history"))
			router.Handler(createRequestCtxFromPath("GET", "/account3/subscription3"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for the index route", func() {
			value1 := 1

			router.GET("/account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/"))

			Expect(value1).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for static routes", func() {
			value1 := 1

			router.GET("/account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/account/transactions_notfound"))

			Expect(value1).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for static routes half path", func() {
			value1 := 1

			router.GET("/account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/account"))

			Expect(value1).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for wildcard routes", func() {
			value1 := 1
			value2 := 1
			value3 := 1

			router.GET("/:account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})
			router.GET("/:account/profile", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})
			router.GET("/:user/roles", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/value1/transactions_notfound"))

			router.NotFound = func(ctx *Context) {
				value2 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/value2/profile_notfound"))

			router.NotFound = func(ctx *Context) {
				value3 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/value3/roles_notfound"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for wildcard half path", func() {
			value1 := 1

			router.GET("/:account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("GET", "/value1"))

			Expect(value1).To(Equal(2))
		})

		ginkgo.It("should call the not found callback for wrong method", func() {
			value1 := 1

			router.GET("/:account/transactions", func(ctx *Context) {
				ginkgo.Fail("should not be called")
			})

			router.NotFound = func(ctx *Context) {
				value1 = 2
			}
			router.Handler(createRequestCtxFromPath("POST", "/value1"))

			Expect(value1).To(Equal(2))
		})

		ginkgo.Describe("Middlewares", func() {
			ginkgo.It("should call all the middlewares in sequence", func() {
				calls := make([]string, 0)
				router.GET("/:account/transactions", func(ctx *Context) {
					calls = append(calls, "endpoint")
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "middleware1")
					next(ctx)
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "middleware2")
					next(ctx)
				})
				router.Handler(createRequestCtxFromPath("GET", "/account/transactions"))
				Expect(calls).To(HaveLen(3))
				Expect(calls[0]).To(Equal("middleware1"))
				Expect(calls[1]).To(Equal("middleware2"))
				Expect(calls[2]).To(Equal("endpoint"))
			})

			ginkgo.It("should the middleware prevent a handler and  for being called", func() {
				calls := make([]string, 0)
				router.GET("/:account/transactions", func(ctx *Context) {
					ginkgo.Fail("this endpoint should not be called")
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "middleware1")
				}, func(ctx *Context, next Handler) {
					ginkgo.Fail("this middleware should not be called")
				})
				router.Handler(createRequestCtxFromPath("GET", "/account/transactions"))
				Expect(calls).To(HaveLen(1))
				Expect(calls[0]).To(Equal("middleware1"))
			})

			ginkgo.It("should call the group middleware for a route", func() {
				calls := make([]string, 0)
				group := router.Group("/v1", func(ctx *Context, next Handler) {
					calls = append(calls, "groupMiddleware1")
					next(ctx)
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "groupMiddleware2")
					next(ctx)
				})
				subgroup := group.Group("/subgroup", func(ctx *Context, next Handler) {
					calls = append(calls, "subgroupMiddleware1")
					next(ctx)
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "subgroupMiddleware2")
					next(ctx)
				})
				subgroup.GET("/route1", func(ctx *Context) {
					calls = append(calls, "endpoint")
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "middleware1")
					next(ctx)
				})
				router.Handler(createRequestCtxFromPath("GET", "/v1/subgroup/route1"))
				Expect(calls).To(HaveLen(6))
				Expect(calls[0]).To(Equal("groupMiddleware1"))
				Expect(calls[1]).To(Equal("groupMiddleware2"))
				Expect(calls[2]).To(Equal("subgroupMiddleware1"))
				Expect(calls[3]).To(Equal("subgroupMiddleware2"))
				Expect(calls[4]).To(Equal("middleware1"))
				Expect(calls[5]).To(Equal("endpoint"))
			})

			ginkgo.It("should call the group middleware avoid calling the next middleware and the route", func() {
				calls := make([]string, 0)
				group := router.Group("/v1", func(ctx *Context, next Handler) {
					calls = append(calls, "groupMiddleware1")
					next(ctx)
				})
				subgroup := group.Group("/subgroup", func(ctx *Context, next Handler) {
					calls = append(calls, "subgroupMiddleware1")
					next(ctx)
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "subgroupMiddleware2")
				})
				subgroup.GET("/route1", func(ctx *Context) {
					calls = append(calls, "endpoint")
				}, func(ctx *Context, next Handler) {
					calls = append(calls, "middleware1")
					next(ctx)
				})
				router.Handler(createRequestCtxFromPath("GET", "/v1/subgroup/route1"))
				Expect(calls).To(HaveLen(3))
				Expect(calls[0]).To(Equal("groupMiddleware1"))
				Expect(calls[1]).To(Equal("subgroupMiddleware1"))
				Expect(calls[2]).To(Equal("subgroupMiddleware2"))
			})
		})
	})
})

func BenchmarkSplit(b *testing.B) {
	path := []byte("/path/with/four/parts")
	tokens := make([][]byte, 0)

	for i := 0; i < b.N; i++ {
		tokens = Split(path, tokens)
		tokens = tokens[0:0]
	}
}

func BenchmarkRouter_Handler(b *testing.B) {
	router := NewRouter()
	router.GET("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	for i := 0; i < b.N; i++ {
		router.Handler(&ctx)
	}
}

func BenchmarkRouter_HandlerWithMiddleware(b *testing.B) {
	router := NewRouter()
	router.GET("/", emptyHandler, func(ctx *Context, next Handler) {
		next(ctx)
	})
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	for i := 0; i < b.N; i++ {
		router.Handler(&ctx)
	}
}
