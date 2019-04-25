package main

import (
	"fmt"

	"github.com/lab259/http"
	"github.com/lab259/http/examples/todos/api"
	"github.com/lab259/http/examples/todos/errors"
)

var config = http.Config{
	Name: "TODO App",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewRouter(notFound)
	api.SetupRoutes(router)
	return router
}

func main() {
	app := http.NewApplication(config, router())
	fmt.Printf("%s listening at %s...\n", app.Name(), app.Configuration.HTTP.Bind)
	app.Start()
}

func notFound(req http.Request, res http.Response) http.Result {
	return res.Status(404).Error(errors.ErrNotFound)
}
