package main

import (
	"fmt"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
)

var config = hermes.ApplicationConfig{
	Name: "cors/v0.1.0",
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.DefaultRouter()
	router.Use(middlewares.DefaultCorsMiddleware())
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
