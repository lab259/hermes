package main

import (
	"fmt"
	"time"

	"github.com/lab259/hermes/middlewares"

	"github.com/lab259/hermes"
)

var config = hermes.ApplicationConfig{
	Name: "Hello World (w/ Middlewares)",
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.NewRouter(hermes.RouterConfig{})
	router.Use(
		middlewares.RecoverableMiddleware,
		logMiddleware,
	)

	router.Get("/hello", func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})

	router.Get("/crash", func(req hermes.Request, res hermes.Response) hermes.Result {
		panic("oops")
	})

	return router
}

func main() {
	app := hermes.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/hello")
	fmt.Println("Go to http://localhost:8080/crash")
	app.Start()
}

func logMiddleware(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
	now := time.Now()
	defer fmt.Printf("%s [%d] %s: %s (took %s)\n", now.UTC().Format(time.RFC3339), req.Raw().ID(), req.Method(), req.Path(), time.Since(now))
	return next(req, res)
}
