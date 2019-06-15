package main

import (
	"fmt"
	"time"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
)

var config = hermes.ApplicationConfig{
	Name: "files/v0.1.0",
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.NewRouter(hermes.RouterConfig{})
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/view", func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.File("examples/files/sample.pdf")
	})
	router.Get("/download", func(req hermes.Request, res hermes.Response) hermes.Result {
		now := time.Now().UTC()
		return res.FileDownload("examples/files/sample.pdf", fmt.Sprintf("sample-%d.pdf", now.Unix()))
	})
	router.Get("/file", func(req hermes.Request, res hermes.Response) hermes.Result {
		qs := hermes.ParseQuery(req)
		if qs.Bool("download") {
			return res.FileDownload("examples/files/sample.pdf", "sample.pdf")
		}
		return res.File("examples/files/sample.pdf")
	})
	return router
}

func main() {
	app := hermes.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/view")
	fmt.Println("Go to http://localhost:8080/download")
	fmt.Println("Go to http://localhost:8080/file")
	app.Start()
}
