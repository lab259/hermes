package main

import (
	"fmt"

	"github.com/lab259/http"
	"github.com/lab259/http/middlewares"
)

var config = http.Config{
	Name: "Hello World",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewRouter(nil)
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/hello", func(req http.Request, res http.Response) http.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})
	return router
}

func main() {
	app := http.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/hello")
	app.Start()
}
