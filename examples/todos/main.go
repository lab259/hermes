package main

import (
	"fmt"

	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/api"
	"github.com/lab259/http/middlewares"
)

var config = http.ApplicationConfig{
	Name: "TODO API",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewDefaultRouter()

	router.Use(
		middlewares.RecoverableMiddleware,
		middlewares.LoggingMiddleware,
	)

	api.SetupRoutes(router)
	return router
}

func main() {
	app := http.NewApplication(config, router())
	fmt.Printf("%s listening at http://localhost%s ...\n", app.Name(), app.Configuration.HTTP.Bind)
	app.Start()
}
