package main

import (
	"fmt"
	"time"

	"github.com/lab259/go-rscsrv"

	"github.com/lab259/http"
	"github.com/lab259/http/middlewares"
)

var config = http.ApplicationConfig{
	Name: "Hello World",
	ServiceStarter: rscsrv.NewServiceStarter([]rscsrv.Service{
		&serviceA{},
		&serviceB{},
	}, &rscsrv.ColorServiceReporter{}),
	HTTP: http.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() http.Router {
	router := http.NewRouter(http.RouterConfig{})
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/hello", func(req http.Request, res http.Response) http.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})
	return router
}

func main() {
	app := http.NewApplication(config, router())
	app.Configuration.ServiceStarter.Start()

	fmt.Println("Go to http://localhost:8080/hello")
	app.Start()
}

/** Service A **/

type serviceA struct {
}

func (service *serviceA) Name() string {
	return "Service A"
}

func (service *serviceA) LoadConfiguration() (interface{}, error) {
	time.Sleep(time.Second)
	return nil, nil
}

func (service *serviceA) ApplyConfiguration(interface{}) error {
	return nil
}

func (service *serviceA) Restart() error {
	if err := service.Stop(); err != nil {
		return err
	}
	return service.Start()
}

func (service *serviceA) Start() error {
	time.Sleep(time.Second)
	return nil
}

func (service *serviceA) Stop() error {
	time.Sleep(time.Second)
	return nil
}

/** ServiceB **/

type serviceB struct {
}

func (service *serviceB) Name() string {
	return "Service B"
}

func (service *serviceB) LoadConfiguration() (interface{}, error) {
	time.Sleep(time.Second)
	return nil, nil
}

func (service *serviceB) ApplyConfiguration(interface{}) error {
	return nil
}

func (service *serviceB) Restart() error {
	if err := service.Stop(); err != nil {
		return err
	}
	return service.Start()
}

func (service *serviceB) Start() error {
	time.Sleep(time.Second)
	return nil
}

func (service *serviceB) Stop() error {
	time.Sleep(time.Second)
	return nil
}
