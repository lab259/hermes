package main

import (
	"fmt"
	"time"

	"github.com/lab259/http"
)

var config = http.Config{
	Name: "Hello World (w/ Middlewares)",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() *http.Router {
	router := http.NewRouter()
	router.Use(
		recoverMiddleware,
		logMiddleware,
	)

	router.GET("/hello", func(ctx *http.Context) {
		ctx.SendJson(map[string]interface{}{
			"hello": "world",
		})
	})

	router.GET("/crash", func(ctx *http.Context) {
		panic("oops")
	})

	return router
}

func main() {
	app := http.NewApplication(config, router())
	app.Start()
}

func recoverMiddleware(ctx *http.Context, next http.Handler) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%s %s [recovered: %s]\n", time.Now().UTC().Format(time.RFC3339), ctx.Request.RequestURI(), e)
			ctx.SendJson(map[string]interface{}{
				"error": 500,
			})
		}
	}()
	next(ctx)
}

func logMiddleware(ctx *http.Context, next http.Handler) {
	now := time.Now()
	next(ctx)
	fmt.Printf("%s %s [%s]\n", now.UTC().Format(time.RFC3339), ctx.Request.RequestURI(), time.Since(now))
}
