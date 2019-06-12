package http

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type ApplicationConfig struct {
	Name string
	HTTP FasthttpServiceConfiguration
}

type Application struct {
	fasthttpService FasthttpService
	Configuration   ApplicationConfig
}

func NewApplication(config ApplicationConfig, router Router) *Application {
	app := &Application{
		Configuration: config,
	}

	if config.Name != "" {
		app.fasthttpService.Server.Name = fmt.Sprintf("fasthttp/%s", config.Name)
	}

	app.fasthttpService.Server.Handler = router.Handler()

	return app
}

func (app *Application) Name() string {
	if app.Configuration.Name == "" {
		return "Application"
	}
	return app.Configuration.Name
}

func (app *Application) LoadConfiguration() (interface{}, error) {
	return nil, nil
}

func (app *Application) ApplyConfiguration(interface{}) error {
	return nil
}

func (app *Application) Restart() error {
	err := app.Stop()
	if err != nil {
		return err
	}
	return app.Start()
}

func (app *Application) Start() error {
	err := app.fasthttpService.ApplyConfiguration(app.Configuration.HTTP)
	if err != nil {
		return err
	}

	done := make(chan bool, 1)
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		app.Stop()
		close(done)
	}()

	if err := app.fasthttpService.Start(); err != nil {
		return err
	}

	<-done
	return nil
}

func (app *Application) Stop() error {
	return app.fasthttpService.Stop()
}
