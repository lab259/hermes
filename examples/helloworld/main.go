package main

import (
	"fmt"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
)

var config = hermes.ApplicationConfig{
	Name: "Hello World",
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.NewRouter(hermes.RouterConfig{})
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/hello", func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})
	return router
}

func main() {
	app := hermes.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/hello")
	app.Start()
}
