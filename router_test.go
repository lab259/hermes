package http

import (
	"strings"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/valyala/fasthttp"
)

func createRequestCtxFromPath(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.Request.URI().SetPath(path)
	return ctx
}

var emptyRouterConfig = RouterConfig{
	NotFound: emptyHandler,
}

var emptyResult = &result{}

var emptyHandler = func(req Request, res Response) Result {
	return res.End()
}

func createPathDescriptor() *tokensDescriptor {
	return &tokensDescriptor{
		m: make(map[int][]byte, 10),
		n: 0,
	}
}

var _ = g.Describe("Router", func() {

	g.Describe("Split", func() {
		g.It("should split the path", func() {
			path := []byte("/path/with/four/parts")

			tokens := createPathDescriptor()
			split(path, tokens)
			Expect(tokens.n).To(Equal(4))
			Expect(tokens.m).To(HaveLen(4))
			Expect(tokens.m[0]).To(Equal([]byte("path")))
			Expect(tokens.m[1]).To(Equal([]byte("with")))
			Expect(tokens.m[2]).To(Equal([]byte("four")))
			Expect(tokens.m[3]).To(Equal([]byte("parts")))
		})

		g.It("should split the path not starting with /", func() {
			path := []byte("path/with/four/parts")
			tokens := createPathDescriptor()
			split(path, tokens)
			Expect(tokens.n).To(Equal(4))
			Expect(tokens.m).To(HaveLen(4))
			Expect(tokens.m[0]).To(Equal([]byte("path")))
			Expect(tokens.m[1]).To(Equal([]byte("with")))
			Expect(tokens.m[2]).To(Equal([]byte("four")))
			Expect(tokens.m[3]).To(Equal([]byte("parts")))
		})

		g.It("should split the path ending with /", func() {
			path := []byte("/path/with/four/parts/")
			tokens := createPathDescriptor()
			split(path, tokens)
			Expect(tokens.n).To(Equal(4))
			Expect(tokens.m).To(HaveLen(4))
			Expect(tokens.m[0]).To(Equal([]byte("path")))
			Expect(tokens.m[1]).To(Equal([]byte("with")))
			Expect(tokens.m[2]).To(Equal([]byte("four")))
			Expect(tokens.m[3]).To(Equal([]byte("parts")))
		})

		g.It("should split an empty path", func() {
			path := []byte("/")
			tokens := createPathDescriptor()
			split(path, tokens)
			Expect(tokens.m).To(BeEmpty())
			Expect(tokens.n).To(Equal(0))
		})

		// g.FIt("should split an empty path", func() {
		// 	path := []byte("/")
		// 	tokens := make([][]byte, 0)
		// 	tokens = bytes.Split(path, []byte{'/'})
		// 	Expect(tokens).To(BeEmpty())
		// })
	})

	g.Describe("Parse", func() {

		g.It("should parse a GET", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(BeEmpty())
			Expect(router.children["GET"].wildcard).To(BeNil())
			Expect(router.children["GET"].handler).NotTo(BeNil())
		})

		g.It("should parse a GET²", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/route", emptyHandler)

			Expect(router.children).To(HaveKey("GET"))
			Expect(router.children["GET"].children).To(HaveKey("route"))
			Expect(router.children["GET"].wildcard).To(BeNil())
		})

		g.It("should parse a POST", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Post("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		g.It("should parse a PUT", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Put("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PUT"))
			Expect(router.children["PUT"].children).To(HaveKey("route"))
			Expect(router.children["PUT"].wildcard).To(BeNil())
		})

		g.It("should parse a DELETE", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Delete("/route", emptyHandler)

			Expect(router.children).To(HaveKey("DELETE"))
			Expect(router.children["DELETE"].children).To(HaveKey("route"))
			Expect(router.children["DELETE"].wildcard).To(BeNil())
		})

		g.It("should parse a HEAD", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Head("/route", emptyHandler)

			Expect(router.children).To(HaveKey("HEAD"))
			Expect(router.children["HEAD"].children).To(HaveKey("route"))
			Expect(router.children["HEAD"].wildcard).To(BeNil())
		})

		g.It("should parse a OPTIONS", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Options("/route", emptyHandler)

			Expect(router.children).To(HaveKey("OPTIONS"))
			Expect(router.children["OPTIONS"].children).To(HaveKey("route"))
			Expect(router.children["OPTIONS"].wildcard).To(BeNil())
		})

		g.It("should parse a PATCH", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Patch("/route", emptyHandler)

			Expect(router.children).To(HaveKey("PATCH"))
			Expect(router.children["PATCH"].children).To(HaveKey("route"))
			Expect(router.children["PATCH"].wildcard).To(BeNil())
		})

		g.It("should parse a POST", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Post("/route", emptyHandler)

			Expect(router.children).To(HaveKey("POST"))
			Expect(router.children["POST"].children).To(HaveKey("route"))
			Expect(router.children["POST"].wildcard).To(BeNil())
		})

		g.It("should parse a complete static route", func() {
			router := NewRouter(emptyRouterConfig).(*router)
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
		})

		g.It("should parse multiple static routes related", func() {
			router := NewRouter(emptyRouterConfig).(*router)
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

			Expect(router.children["GET"].children["this"].children).To(HaveKey("should2"))
			Expect(router.children["GET"].children["this"].children["should2"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].children).To(HaveKey("static"))
		})

		g.It("should parse a complete a route starting static and ending with a wildcard", func() {
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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

			Expect(router.children["GET"].children["this"].children).To(HaveKey("should2"))
			Expect(router.children["GET"].children["this"].children["should2"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children).To(HaveKey("be"))
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].wildcard).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].handler).To(BeNil())
			Expect(router.children["GET"].children["this"].children["should2"].children["be"].children).To(HaveKey("static"))

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
		})

		g.It("should parse a complete route with wildcard", func() {
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)

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
			router := NewRouter(emptyRouterConfig).(*router)

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
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/account/detail", emptyHandler)
			Expect(func() {
				router.Get("/account/detail", emptyHandler)
			}).To(Panic())
		})

		g.It("should panic due to conflicting 'wildcarded' routes", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/:account", emptyHandler)
			Expect(func() {
				router.Get("/:transaction", emptyHandler)
			}).To(Panic())
		})

		g.It("should panic due to conflicting mixing routes", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/:account/detail", emptyHandler)
			router.Get("/:account/id", emptyHandler)
			Expect(func() {
				router.Get("/:transaction/id", emptyHandler)
			}).To(Panic())
		})

		g.It("should not match any ropute", func() {
			router := NewRouter(emptyRouterConfig).(*router)
			router.Get("/:account/detail", emptyHandler)
			router.Get("/:account/id", emptyHandler)
			ok, _ := router.children["GET"].Matches(0, createPathDescriptor(), nil)
			Expect(ok).To(BeFalse())
		})
	})

	g.Describe("Group", func() {
		g.It("should parse a GET", func() {
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
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
			router := NewRouter(emptyRouterConfig).(*router)
			group := router.Prefix("/group").(*route)
			group2 := group.Prefix("/subgroup").(*route)

			Expect(group.prefix).To(Equal("group"))
			Expect(group2).NotTo(BeNil())
			Expect(group2.prefix).To(Equal("group/subgroup"))
		})
	})

	g.Describe("Handle", func() {
		var router Router

		g.BeforeEach(func() {
			router = NewRouter(emptyRouterConfig)
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

			router := NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			}})

			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for static routes", func() {
			value1 := 1

			router = NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			}})
			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			router.Handler()(createRequestCtxFromPath("GET", "/account/transactions_notfound"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for static routes half path", func() {
			value1 := 1

			router := NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			}})
			router.Get("/account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/account"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the not found callback for wildcard routes", func() {
			value := 0
			router := NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
				value++
				return res.End()
			}})
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
			router := NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
				value1 = 2
				return res.End()
			}})
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			router.Handler()(createRequestCtxFromPath("GET", "/value1"))

			Expect(value1).To(Equal(2))
		})

		g.It("should call the default not found callback", func() {
			router := NewDefaultRouter()
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			ctx := createRequestCtxFromPath("GET", "/value1")
			router.Handler()(ctx)
			Expect(ctx.Response.StatusCode()).To(Equal(fasthttp.StatusNotFound))
		})

		g.It("should call the default method not allowed handler for wrong method", func() {
			router := NewDefaultRouter()
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			ctx := createRequestCtxFromPath("POST", "/value1/transactions")
			router.Handler()(ctx)
			methods := string(ctx.Response.Header.Peek("Allow"))
			Expect(strings.Split(methods, ", ")).To(ConsistOf("GET", "OPTIONS"))
			Expect(ctx.Response.StatusCode()).To(Equal(fasthttp.StatusMethodNotAllowed))
		})

		g.It("should call the default not found callback (JSON)", func() {
			router := NewDefaultRouter()
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})
			ctx := createRequestCtxFromPath("GET", "/value1")
			ctx.Request.Header.Set("Accept", "application/json, text/html, text/plain")
			router.Handler()(ctx)
			Expect(ctx.Response.StatusCode()).To(Equal(fasthttp.StatusNotFound))
			Expect(string(ctx.Response.Header.Peek("Content-Type"))).To(Equal("application/json; charset=utf-8"))
		})

		g.It("should call the default method not allowed handler for wrong method (JSON)", func() {
			router := NewDefaultRouter()
			router.Get("/:account/transactions", func(req Request, res Response) Result {
				g.Fail("should not be called")
				return res.End()
			})

			ctx := createRequestCtxFromPath("POST", "/value1/transactions")
			ctx.Request.Header.Set("Accept", "application/json, text/html, text/plain")
			router.Handler()(ctx)
			methods := string(ctx.Response.Header.Peek("Allow"))
			Expect(strings.Split(methods, ", ")).To(ConsistOf("GET", "OPTIONS"))
			Expect(ctx.Response.StatusCode()).To(Equal(fasthttp.StatusMethodNotAllowed))
			Expect(string(ctx.Response.Header.Peek("Content-Type"))).To(Equal("application/json; charset=utf-8"))
		})

		g.Describe("Options request", func() {
			g.It("should resolve server-wide", func() {
				router.Get("/todos", emptyHandler)
				router.Post("/todos", emptyHandler)
				ctx := createRequestCtxFromPath("OPTIONS", "*")
				router.Handler()(ctx)
				methods := string(ctx.Response.Header.Peek("Allow"))
				Expect(strings.Split(methods, ", ")).To(ConsistOf("GET", "POST", "OPTIONS"))
			})

			g.It("should resolve server-wide²", func() {
				router.Get("/todos", emptyHandler)
				router.Post("/todos", emptyHandler)
				ctx := createRequestCtxFromPath("OPTIONS", "/*")
				router.Handler()(ctx)
				methods := string(ctx.Response.Header.Peek("Allow"))
				Expect(strings.Split(methods, ", ")).To(ConsistOf("GET", "POST", "OPTIONS"))
			})

			g.It("should resolve specific path", func() {
				router.Get("/todos", emptyHandler)
				router.Post("/todos", emptyHandler)
				ctx := createRequestCtxFromPath("OPTIONS", "/todos")
				router.Handler()(ctx)
				methods := string(ctx.Response.Header.Peek("Allow"))
				Expect(strings.Split(methods, ", ")).To(ConsistOf("GET", "POST", "OPTIONS"))
			})
		})

		g.Describe("Middlewares", func() {
			g.It("should call root middlewares with not found", func() {
				calls := make([]string, 0)
				router := NewRouter(RouterConfig{NotFound: func(req Request, res Response) Result {
					calls = append(calls, "notfound")
					return res.End()
				}})
				router.Use(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware1")
					return next(req, res)
				}, func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware2")
					return next(req, res)
				})
				router.Get("/:account/transactions", func(req Request, res Response) Result {
					calls = append(calls, "endpoint")
					return res.End()
				})
				router.Handler()(createRequestCtxFromPath("GET", "/account_not_found"))
				Expect(calls).To(HaveLen(3))
				Expect(calls[0]).To(Equal("middleware1"))
				Expect(calls[1]).To(Equal("middleware2"))
				Expect(calls[2]).To(Equal("notfound"))
			})

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

			g.It("should not call middlewares added after prefix/group", func() {
				calls := make([]string, 0)
				router.Use(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware1")
					return next(req, res)
				}, func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware2")
					return next(req, res)
				})

				router.Prefix("/api").Get("/:account/transactions", func(req Request, res Response) Result {
					calls = append(calls, "endpoint")
					return res.End()
				})

				router.Use(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "middleware3")
					return next(req, res)
				})
				router.Handler()(createRequestCtxFromPath("GET", "/api/account/transactions"))
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

				group.Use(func(req Request, res Response, next Handler) Result {
					calls = append(calls, "groupMiddleware3")
					return next(req, res)
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
	tokens := createPathDescriptor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		split(path, tokens)
		tokens.n = 0
	}
}

func BenchmarkRouter_Handler(b *testing.B) {
	router := NewRouter(RouterConfig{})
	router.Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}

func BenchmarkRouter_HandlerWithParams(b *testing.B) {
	router := NewRouter(RouterConfig{})
	router.Get("/:id", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/id")

	h := router.Handler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}

func BenchmarkRouter_HandlerWithMiddleware(b *testing.B) {
	router := NewRouter(RouterConfig{})
	router.With(func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}).Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}

func BenchmarkRouter_HandlerWithMiddleware2(b *testing.B) {
	router := NewRouter(RouterConfig{})
	router.With(func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}, func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}).Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}

func BenchmarkRouter_HandlerWithMiddleware5(b *testing.B) {
	router := NewRouter(RouterConfig{})
	router.With(func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}, func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}, func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}, func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}, func(req Request, res Response, next Handler) Result {
		return next(req, res)
	}).Get("/", emptyHandler)
	ctx := fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/")

	h := router.Handler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(&ctx)
	}
}
