package main

import (
	"fmt"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/examples/todos/api"
	"github.com/lab259/hermes/middlewares"
)

var config = hermes.ApplicationConfig{
	Name: "Todo (v0.1.0)",
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.DefaultRouter()

	router.Use(
		middlewares.RecoverableMiddleware,
		middlewares.LoggingMiddleware,
	)

	api.SetupRoutes(router)
	return router
}

func main() {
	app := hermes.NewApplication(config, router())
	fmt.Printf("%s listening at http://localhost%s ...\n", app.Name(), app.Configuration.HTTP.Bind)
	app.Start()
}
