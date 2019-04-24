package http

type Config struct {
	Name string
	HTTP FasthttpServiceConfiguration
}

type Application struct {
	fasthttpService FasthttpService
	Configuration   Config
}

func (service *Application) Name() string {
	return service.Configuration.Name
}

func NewApplication(config Config, router Router) *Application {
	app := &Application{
		Configuration: config,
	}

	app.fasthttpService.Server.Handler = router.Handler()
	return app
}

func (app *Application) Start() error {
	err := app.fasthttpService.ApplyConfiguration(app.Configuration.HTTP)
	if err != nil {
		return err
	}
	return app.fasthttpService.Start()
}
