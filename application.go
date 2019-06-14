package http

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lab259/go-rscsrv"
)

type ApplicationConfig struct {
	Name           string
	ServiceStarter *rscsrv.ServiceStarter
	HTTP           FasthttpServiceConfiguration
}

type Application struct {
	fasthttpService FasthttpService
	Configuration   ApplicationConfig
	done            chan bool
	signals         chan os.Signal
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
	if err := app.Stop(); err != nil {
		return err
	}
	return app.Start()
}

func (app *Application) Start() error {
	err := app.fasthttpService.ApplyConfiguration(app.Configuration.HTTP)
	if err != nil {
		return err
	}

	app.done = make(chan bool, 1)
	go func() {
		app.signals = make(chan os.Signal, 1)
		signal.Notify(app.signals, syscall.SIGINT, syscall.SIGTERM)
		if _, ok := <-app.signals; ok {
			app.Stop()
		}
	}()

	if err := app.fasthttpService.Start(); err != nil {
		return err
	}

	<-app.done
	return nil
}

func (app *Application) Stop() error {
	defer func() {
		if app.Configuration.ServiceStarter != nil {
			app.Configuration.ServiceStarter.Stop(true)
		}
		close(app.signals)
		close(app.done)
	}()
	return app.fasthttpService.Stop()
}
