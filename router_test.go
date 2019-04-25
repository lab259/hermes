package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

func createRequestCtxFromPath(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.URI().SetPath(path)
	return ctx
}

var emptyHandler = func(req Request, res Response) Result {
	return res.End()
}

var _ = g.Describe("Router", func() {

	g.Describe("Split", func() {
		g.It("should split the path", func() {
			path := []byte("/path/with/four/parts")
			tokens := make([][]byte, 0)
			tokens = split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		g.It("should split the path not starting with /", func() {
			path := []byte("path/with/four/parts")
			tokens := make([][]byte, 0)
			tokens = split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		g.It("should split the path ending with /", func() {
			path := []byte("/path/with/four/parts/")
			tokens := make([][]byte, 0)
			tokens = split(path, tokens)
			Expect(tokens).To(HaveLen(4))
			Expect(tokens[0]).To(Equal([]byte("path")))
			Expect(tokens[1]).To(Equal([]byte("with")))
			Expect(tokens[2]).To(Equal([]byte("four")))
			Expect(tokens[3]).To(Equal([]byte("parts")))
		})

		g.It("should split an empty path", func() {
			path := []byte("/")
			tokens := make([][]byte, 0)
			tokens = split(path, tokens)
			Expect(tokens).To(BeEmpty())
		})
	})

	g.Describe("Parse", func() {

		g.It("should parse a GET", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].handler).NotTo(BeNil())
		})

		g.It("should parse a GET", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/route", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("route"))
			Expect(router.children["GET"].wildcard).To(BeNil())
		})

		g.It("should parse a POST", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Post("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		g.It("should parse a PUT", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Put("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PUT"))
			Expect(router.children["PUT"].children).To(HaveKey("route"))
			Expect(router.children["PUT"].wildcard).To(BeNil())
		})

		g.It("should parse a DELETE", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Delete("/route", emptyHandler)

			Expect(router.children).To(HaveKey("DELETE"))
			Expect(router.children["DELETE"].children).To(HaveKey("route"))
			Expect(router.children["DELETE"].wildcard).To(BeNil())
		})

		g.It("should parse a HEAD", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Head("/route", emptyHandler)

			Expect(router.children).To(HaveKey("HEAD"))
			Expect(router.children["HEAD"].children).To(HaveKey("route"))
			Expect(router.children["HEAD"].wildcard).To(BeNil())
		})

		g.It("should parse a OPTIONS", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Options("/route", emptyHandler)

			Expect(router.children).To(HaveKey("OPTIONS"))
			Expect(router.children["OPTIONS"].children).To(HaveKey("route"))
			Expect(router.children["OPTIONS"].wildcard).To(BeNil())
		})

		g.It("should parse a PATCH", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Patch("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PATCH"))
			Expect(router.children["PATCH"].children).To(HaveKey("route"))
			Expect(router.children["PATCH"].wildcard).To(BeNil())
		})

		g.It("should parse a POST", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Post("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		g.It("should parse a complete static route", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/this/should/be/static", emptyHandler)

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

		g.It("should parse multiple static routes related", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/this/should/be/static", emptyHandler)
			router.Get("/this/should2/be/static", emptyHandler)

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

		g.It("should parse a complete a route starting static and ending with a wildcard", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/static/:wildcard", emptyHandler)

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

		g.It("should parse multiple static routes related and not", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/this/should/be/static", emptyHandler)
			router.Get("/this/should2/be/static", emptyHandler)
			router.Get("/this2/should/be/static", emptyHandler)

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

		g.It("should parse a complete route with wildcard", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account/detail/another", emptyHandler)

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

		g.It("should parse a complete route with a sequence of wildcards", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account/:transaction/:invoice", emptyHandler)

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

		g.It("should parse multiple routes starting with wildcards", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account/detail", emptyHandler)
			router.Get("/:account/history", emptyHandler)
			router.Get("/:transaction/invoice", emptyHandler)

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

		g.It("should parse multiple mixed routes", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/accounts/:account/detail", emptyHandler)
			router.Get("/accounts/:account/history", emptyHandler)
			router.Get("/:transaction/invoice", emptyHandler)

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

		g.It("should panic due to conflicting empty tokens", func() {
			router := NewRouter(emptyHandler).(*router)

			Expect(func() {
				router.Get("//detail", emptyHandler)
			}).To(Panic())

			Expect(func() {
				router.Get("/account/detail//", emptyHandler)
			}).To(Panic())

			Expect(func() {
				router.Get("/account//detail", emptyHandler)
			}).To(Panic())
		})

		g.It("should not panic with empty token at the end", func() {
			router := NewRouter(emptyHandler).(*router)

			Expect(func() {
				router.Get("/account/", emptyHandler)
			}).NotTo(Panic())

			Expect(func() {
				router.Get("/account/detail/", emptyHandler)
			}).NotTo(Panic())

			Expect(func() {
				router.Get("/account/detail/:id/", emptyHandler)
			}).NotTo(Panic())
		})

		g.It("should panic due to conflicting static routes", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/account/detail", emptyHandler)
			Expect(func() {
				router.Get("/account/detail", emptyHandler)
			}).To(Panic())
		})

		g.It("should panic due to conflicting 'wildcarded' routes", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account", emptyHandler)
			Expect(func() {
				router.Get("/:transaction", emptyHandler)
			}).To(Panic())
		})

		g.It("should panic due to conflicting mixing routes", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account/detail", emptyHandler)
			router.Get("/:account/id", emptyHandler)
			Expect(func() {
				router.Get("/:transaction/id", emptyHandler)
			}).To(Panic())
		})

		g.It("should not match any ropute", func() {
			router := NewRouter(emptyHandler).(*router)
			router.Get("/:account/detail", emptyHandler)
			router.Get("/:account/id", emptyHandler)
			ok, _, _ := router.children["GET"].Matches(nil, nil)
			Expect(ok).To(BeFalse())
		})
	})

	g.Describe("Group", func() {
		g.It("should parse a GET", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Get("/route", emptyHandler)

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

		g.It("should parse a POST", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Post("/route", emptyHandler)

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

		g.It("should parse a PUT", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Put("/route", emptyHandler)

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

		g.It("should parse a DELETE", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Delete("/route", emptyHandler)

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

		g.It("should parse a HEAD", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Head("/route", emptyHandler)

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

		g.It("should parse a OPTIONS", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Options("/route", emptyHandler)

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

		g.It("should parse a PATCH", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group")
			group.Patch("/route", emptyHandler)

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

		g.It("should check the subgroup", func() {
			router := NewRouter(emptyHandler).(*router)
			group := router.Prefix("/group").(*route)
			group2 := group.Prefix("/subgroup").(*route)

			Expect(group.prefix).To(Equal("/group"))
			Expect(group2).NotTo(BeNil())
			Expect(group2.prefix).To(Equal("/group/subgroup"))
		})
	})

	g.Describe("Handle", func() {
		var router Router

		g.BeforeEach(func() {
			router = NewRouter(emptyHandler)
		})

		g.It("should resolve an empty route", func() {
			value := 1
			router.Get("", func(req Request, res Response) Result {
				value = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/"))

			Expect(value).To(Equal(2))
		})

		g.It("should resolve an empty trailing route", func() {
			value := 1
			router.Get("/", func(req Request, res Response) Result {
				value = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/"))

			Expect(value).To(Equal(2))
		})

		g.It("should resolve a static route", func() {
			value := 1
			router.Get("/static", func(req Request, res Response) Result {
				value = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/static"))

			Expect(value).To(Equal(2))
		})

		g.It("should resolve a static route not starting with /", func() {
			value := 1
			router.Get("static", func(req Request, res Response) Result {
				value = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/static"))

			Expect(value).To(Equal(2))
		})

		g.It("should resolve multiple static routes", func() {
			value1 := 1
			value2 := 1
			value3 := 1

			router.Get("/static", func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})

			router.Get("/static/second", func(req Request, res Response) Result {
				value2 = 2
				return res.End()
			})

			router.Get("/another", func(req Request, res Response) Result {
				value3 = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/static"))
			router.Handler()(createRequestCtxFromPath("GET", "/static/second"))
			router.Handler()(createRequestCtxFromPath("GET", "/another"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		g.It("should resolve a wildcard route", func() {
			value := 1
			router.Get("/:wildcard", func(req Request, res Response) Result {
				Expect(req.Param("wildcard")).To(Equal("value"))
				value = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/value"))

			Expect(value).To(Equal(2))
		})

		g.It("should resolve a multiple wildcard routes", func() {
			value1 := 1
			value2 := 1
			value3 := 1
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				Expect(req.Param("account")).To(Equal("value1"))
				value1 = 2
				return res.End()
			})
			router.Get("/:account/profile", func(req Request, res Response) Result {
				Expect(req.Param("account")).To(Equal("value2"))
				value2 = 2
				return res.End()
			})
			router.Get("/:user/roles", func(req Request, res Response) Result {
				Expect(req.Param("user")).To(Equal("value3"))
				value3 = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/value1/transactions"))
			router.Handler()(createRequestCtxFromPath("GET", "/value2/profile"))
			router.Handler()(createRequestCtxFromPath("GET", "/value3/roles"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		g.It("should resolve a multiple wildcard routes in sequence", func() {
			value1 := 1
			value2 := 1
			value3 := 1
			router.Get("/:account/:subscription/cancel", func(req Request, res Response) Result {
				Expect(req.Param("account")).To(Equal("account1"))
				Expect(req.Param("subscription")).To(Equal("subscription1"))
				value1 = 2
				return res.End()
			})
			router.Get("/:account/:subscription/history", func(req Request, res Response) Result {
				Expect(req.Param("account")).To(Equal("account2"))
				Expect(req.Param("subscription")).To(Equal("subscription2"))
				value2 = 2
				return res.End()
			})
			router.Get("/:account/:subscription", func(req Request, res Response) Result {
				Expect(req.Param("account")).To(Equal("account3"))
				Expect(req.Param("subscription")).To(Equal("subscription3"))
				value3 = 2
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/account1/subscription1/cancel"))
			router.Handler()(createRequestCtxFromPath("GET", "/account2/subscription2/history"))
			router.Handler()(createRequestCtxFromPath("GET", "/account3/subscription3"))

			Expect(value1).To(Equal(2))
			Expect(value2).To(Equal(2))
			Expect(value3).To(Equal(2))
		})

		g.It("should call the not found callback for the index route", func() {
			value1 := 1

			router := NewRouter(func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})

			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for static routes", func() {
			value1 := 1

			router = NewRouter(func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})
			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			router.Handler()(createRequestCtxFromPath("GET", "/account/transactions_notfound"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for static routes half path", func() {
			value1 := 1

			router := NewRouter(func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})
			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/account"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for wildcard routes", func() {
			value := 0
			router := NewRouter(func(req Request, res Response) Result {
				value++
				return res.End()
			})
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			router.Get("/:account/profile", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			router.Get("/:user/roles", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			handler := router.Handler()

			handler(createRequestCtxFromPath("GET", "/value2/profile_notfound"))
			handler(createRequestCtxFromPath("GET", "/value2/profile_notfound"))
			handler(createRequestCtxFromPath("GET", "/value3/roles_notfound"))

			Expect(value).To(Equal(3))
		})

		g.It("should call the not found callback for wildcard half path", func() {
			value1 := 1
			router := NewRouter(func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/value1"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for wrong method", func() {
			value1 := 1

			router := NewRouter(func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			})
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("POST", "/value1"))

			Expect(value1).To(Equal(2))
		})

		g.Describe("Middlewares", func() {
			g.It("should call all the middlewares in sequence", func() {
				calls := make([]string, 0)
				router.With(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware1")
					return next(req, res)
				}, func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware2")
					return next(req, res)
				}).Get("/:account/transactions", func(req Request, res Response) Result {
					calls = append(calls, "endpoint")
					return res.End()
				})
				router.Handler()(createRequestCtxFromPath("GET", "/account/transactions"))
				Expect(calls).To(HaveLen(3))
				Expect(calls[0]).To(Equal("middleware1"))
				Expect(calls[1]).To(Equal("middleware2"))
				Expect(calls[2]).To(Equal("endpoint"))
			})

			g.It("should the middleware prevent a handler and  for being called", func() {
				calls := make([]string, 0)
				router.With(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware1")
					return res.End()
				}, func(req Request, res Response, next Handler) Result {
					g.Fail("this middleware should not be called")
					return next(req, res)
				}).Get("/:account/transactions", func(req Request, res Response) Result {
					g.Fail("this endpoint should not be called")
					return res.End()
				})
				router.Handler()(createRequestCtxFromPath("GET", "/account/transactions"))
				Expect(calls).To(HaveLen(1))
				Expect(calls[0]).To(Equal("middleware1"))
			})

			g.It("should call the group middleware for a route", func() {
				calls := make([]string, 0)
				group := router.Prefix("/v1")
				group.Use(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "groupMiddleware1")
					return next(req, res)
				}, func(req Request, res Response, next Handler) Result {
					calls = append(calls, "groupMiddleware2")
					return next(req, res)
				})

				group.Prefix("/subgroup").Group(func(r Routable) {
					r.Use(func(req Request, res Response, next Handler) Result {
						calls = append(calls, "subgroupMiddleware1")
						return next(req, res)
					}, func(req Request, res Response, next Handler) Result {
						calls = append(calls, "subgroupMiddleware2")
						return next(req, res)
					})

					r.With(func(req Request, res Response, next Handler) Result {
						calls = append(calls, "middleware1")
						return next(req, res)
					}).Get("/route1", func(req Request, res Response) Result {
						calls = append(calls, "endpoint")
						return res.End()
					})
				})
				router.Handler()(createRequestCtxFromPath("GET", "/v1/subgroup/route1"))
				Expect(calls).To(HaveLen(6))
				Expect(calls[0]).To(Equal("groupMiddleware1"))
				Expect(calls[1]).To(Equal("groupMiddleware2"))
				Expect(calls[2]).To(Equal("subgroupMiddleware1"))
				Expect(calls[3]).To(Equal("subgroupMiddleware2"))
				Expect(calls[4]).To(Equal("middleware1"))
				Expect(calls[5]).To(Equal("endpoint"))
			})

			g.It("should call the group middleware avoid calling the next middleware and the route", func() {
				calls := make([]string, 0)

				router.Prefix("/v1").Group(func(r Routable) {
					r.Use(func(req Request, res Response, next Handler) Result {
						calls = append(calls, "groupMiddleware1")
						return next(req, res)
					})

					subgroup := r.Prefix("/subgroup").With(func(req Request, res Response, next Handler) Result {
						calls = append(calls, "subgroupMiddleware1")
						return next(req, res)
					}, func(req Request, res Response, next Handler) Result {
						calls = append(calls, "subgroupMiddleware2")
						return res.End()
					})

					subgroup.With(func(req Request, res Response, next Handler) Result {
						calls = append(calls, "middleware1")
						return next(req, res)
					}).Get("/route1", func(req Request, res Response) Result {
						calls = append(calls, "endpoint")
						return res.End()
					})
				})
				router.Handler()(createRequestCtxFromPath("GET", "/v1/subgroup/route1"))
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
		tokens = split(path, tokens)
		tokens = tokens[0:0]
	}
}

func BenchmarkRouter_Handler(b *testing.B) {
	router := NewRouter(emptyHandler)
	router.Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}

func BenchmarkRouter_HandlerWithMiddleware(b *testing.B) {
	router := NewRouter(emptyHandler)
	router.With(func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}).Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}
