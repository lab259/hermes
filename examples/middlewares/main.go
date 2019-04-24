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

	router.GET("/hello", func(req http.Request, res http.Response) http.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})

	router.GET("/crash", func(req http.Request, res http.Response) http.Result {
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
		e := recover()
		if e != nil {
			fmt.Printf("%s %s [recovered: %s]\n", time.Now().UTC().Format(time.RFC3339), req.Path(), e)
			if err, ok := e.(error); ok {
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
	defer fmt.Printf("%s %s [%s]\n", now.UTC().Format(time.RFC3339), req.Path(), time.Since(now))
	return next(req, res)
}
