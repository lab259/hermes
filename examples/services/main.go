package main

import (
	"fmt"
	"time"

	"github.com/lab259/go-rscsrv"

	"github.com/lab259/hermes"
	"github.com/lab259/hermes/middlewares"
)

var config = hermes.ApplicationConfig{
	Name: "Hello World",
	ServiceStarter: rscsrv.DefaultServiceStarter(
		&serviceA{},
		&serviceB{},
	),
	HTTP: hermes.FasthttpServiceConfiguration{
		Bind: ":8080",
	},
}

func router() hermes.Router {
	router := hermes.NewRouter(hermes.RouterConfig{})
	router.Use(middlewares.LoggingMiddleware)
	router.Get("/hello", func(req hermes.Request, res hermes.Response) hermes.Result {
		return res.Data(map[string]interface{}{
			"hello": "world",
		})
	})
	return router
}

func main() {
	app := hermes.NewApplication(config, router())
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
