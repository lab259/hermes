package main

import (
	"github.com/lab259/http"
)

var config = http.Config{
	Name: "Hello World",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() *http.Router {
	router := http.NewRouter()
	router.GET("/hello", func(ctx *http.Context) {
		ctx.SendJson(map[string]interface{}{
			"hello": "world",
		})
	})
	return router
}

func main() {
	app := http.NewApplication(config, router())
	app.Start()
}
