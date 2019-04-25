package main

import (
	"fmt"

	validator_v9 "gopkg.in/go-playground/validator.v9"

	"github.com/lab259/errors"
	"github.com/lab259/http"
)

var ErrModule = errors.Module("main")
var ErrNotImplemented = errors.Wrap(errors.New("not implemented"), errors.Http(400), ErrModule, errors.Code("not-implemented"), errors.Message("This endpoint still under construction."))

var config = http.Config{
	Name: "Errors",
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewRouter(nil)
	router.Get("/hello", func(req http.Request, res http.Response) http.Result {
		return res.Error(ErrNotImplemented)
	})
	router.Get("/validation", func(req http.Request, res http.Response) http.Result {
		validator := validator_v9.New()
		model := &Model{}

		if err := validator.Struct(model); err != nil {
			return res.Error(errors.Wrap(err, errors.Validation(), ErrModule))
		}
		return res.Data(model)
	})
	return router
}

func main() {
	app := http.NewApplication(config, router())
	fmt.Println("Go to http://localhost:8080/hello")
	fmt.Println("Go to http://localhost:8080/validation")
	app.Start()
}

type Model struct {
	Name string `validate:"required"`
}
