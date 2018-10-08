package http

type ServiceStarterReporter interface {
	BeforeBegin(service Service)

	BeforeLoadConfiguration(service Service)
	AfterLoadConfiguration(service Service, conf interface{}, err error)

	BeforeApplyConfiguration(service Service)
	AfterApplyConfiguration(service Service, conf interface{}, err error)

	BeforeStart(service Service)
	AfterStart(service Service, err error)
}

type NopServiceReporter struct{}

func (*NopServiceReporter) BeforeBegin(service Service) {}

func (*NopServiceReporter) BeforeLoadConfiguration(service Service) {}

func (*NopServiceReporter) AfterLoadConfiguration(service Service, conf interface{}, err error) {}

func (*NopServiceReporter) BeforeApplyConfiguration(service Service) {}

func (*NopServiceReporter) AfterApplyConfiguration(service Service, conf interface{}, err error) {}

func (*NopServiceReporter) BeforeStart(service Service) {}

func (*NopServiceReporter) AfterStart(service Service, err error) {}

type ServiceStarter struct {
	services []Service
	reporter ServiceStarterReporter
}

func NewServiceStarter(services []Service, reporter ServiceStarterReporter) *ServiceStarter {
	return &ServiceStarter{
		services: services,
		reporter: reporter,
	}
}

func (engineStarter *ServiceStarter) Start() error {
	for _, srv := range engineStarter.services {
		engineStarter.reporter.BeforeBegin(srv)
		engineStarter.reporter.BeforeLoadConfiguration(srv)
		conf, err := srv.LoadConfiguration()
		engineStarter.reporter.AfterLoadConfiguration(srv, conf, err)
		if err != nil {
			return err
		}

		engineStarter.reporter.BeforeApplyConfiguration(srv)
		err = srv.ApplyConfiguration(conf)
		if err != nil {
			return err
		}
		engineStarter.reporter.AfterApplyConfiguration(srv, conf, err)

		engineStarter.reporter.BeforeStart(srv)
		err = srv.Start()
		engineStarter.reporter.AfterStart(srv, err)
		if err != nil {
			return err
		}
	}
	return nil
}
