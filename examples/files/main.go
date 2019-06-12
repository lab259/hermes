package main

import (
	"fmt"
	"time"

	"github.com/lab259/http"
	"github.com/lab259/http/middlewares"
)

var config = http.ApplicationConfig{
	Name: "files/v0.1.0",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewRouter(http.RouterConfig{})
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/view", func(req http.Request, res http.Response) http.Result {
		return res.File("examples/files/sample.pdf")
	})
	router.Get("/download", func(req http.Request, res http.Response) http.Result {
		now := time.Now().UTC()
		return res.FileDownload("examples/files/sample.pdf", fmt.Sprintf("sample-%d.pdf", now.Unix()))
	})
	router.Get("/file", func(req http.Request, res http.Response) http.Result {
		qs := http.ParseQuery(req)
		if qs.Bool("download") {
			return res.FileDownload("examples/files/sample.pdf", "sample.pdf")
		}
		return res.File("examples/files/sample.pdf")
	})
	return router
}

func main() {
	app := http.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/view")
	fmt.Println("Go to http://localhost:8080/download")
	fmt.Println("Go to http://localhost:8080/file")
	app.Start()
}
