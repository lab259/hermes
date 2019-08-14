package hermes

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lab259/go-rscsrv"
)

type ApplicationConfig struct {
	Name           string
	ServiceStarter rscsrv.ServiceStarter
	HTTP           FasthttpServiceConfiguration
}

type Application struct {
	serviceState
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
	app.done = make(chan bool, 1)
	app.signals = make(chan os.Signal, 1)

	return app
}

func (app *Application) Name() string {
	if app.Configuration.Name == "" {
		return "Application"
	}
	return app.Configuration.Name
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

	go func() {
		signal.Notify(app.signals, syscall.SIGINT, syscall.SIGTERM)
		if _, ok := <-app.signals; ok {
			app.Stop()
		}
	}()

	app.setRunning(true)
	if err := app.fasthttpService.Start(); err != nil {
		return err
	}

	<-app.done
	return nil
}

func (app *Application) Stop() error {
	if app.isRunning() {
		defer func() {
			if app.Configuration.ServiceStarter != nil {
				app.Configuration.ServiceStarter.Stop(true)
			}
			signal.Stop(app.signals)
			app.done <- true
			app.setRunning(false)
		}()
		return app.fasthttpService.Stop()
	}
	return nil
}
