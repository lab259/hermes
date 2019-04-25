package main

import (
	"errors"
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

func router() http.Router {
	router := http.NewRouter(nil)
	router.Use(
		recoverMiddleware,
		logMiddleware,
	)

	router.Get("/hello", func(req http.Request, res http.Response) http.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})

	router.Get("/crash", func(req http.Request, res http.Response) http.Result {
		panic("oops")
	})

	return router
}

func main() {
	app := http.NewApplication(config, router())
	app.Start()
}

func recoverMiddleware(req http.Request, res http.Response, next http.Handler) http.Result {
	defer func() {
		recoveredData := recover()
		if recoveredData != nil {
			fmt.Printf("%s [%d] recovered: %s\n", time.Now().UTC().Format(time.RFC3339), req.Raw().ID(), recoveredData)
			if err, ok := recoveredData.(error); ok {
				res.Error(err)
			} else {
				res.Error(errors.New("internal server error"))
			}
		}
	}()
	return next(req, res)
}

func logMiddleware(req http.Request, res http.Response, next http.Handler) http.Result {
	now := time.Now()
	defer fmt.Printf("%s [%d] %s: %s (took %s)\n", now.UTC().Format(time.RFC3339), req.Raw().ID(), req.Method(), req.Path(), time.Since(now))
	return next(req, res)
}
