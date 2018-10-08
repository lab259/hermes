package http_test

import (
	"."
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type countEngineReporter struct {
	countBeforeBegin              int
	countBeforeLoadConfiguration  int
	countAfterLoadConfiguration   int
	countBeforeApplyConfiguration int
	countAfterApplyConfiguration  int
	countBeforeStart              int
	countAfterStart               int
}

func (reporter *countEngineReporter) BeforeBegin(service http.Service) {
	reporter.countBeforeBegin++
}

func (reporter *countEngineReporter) BeforeLoadConfiguration(service http.Service) {
	reporter.countBeforeLoadConfiguration++
}

func (reporter *countEngineReporter) AfterLoadConfiguration(service http.Service, conf interface{}, err error) {
	reporter.countAfterLoadConfiguration++
}

func (reporter *countEngineReporter) BeforeApplyConfiguration(service http.Service) {
	reporter.countBeforeApplyConfiguration++
}

func (reporter *countEngineReporter) AfterApplyConfiguration(service http.Service, conf interface{}, err error) {
	reporter.countAfterApplyConfiguration++
}

func (reporter *countEngineReporter) BeforeStart(service http.Service) {
	reporter.countBeforeStart++
}

func (reporter *countEngineReporter) AfterStart(service http.Service, err error) {
	reporter.countAfterStart++
}

type MockService struct {
	errLoadingConfiguration error
	errApplyConfiguration   error
	errRestart              error
	errStart                error
	errStop                 error
}

func (service *MockService) Name() string {
	return "mock-service"
}

func (service *MockService) LoadConfiguration() (interface{}, error) {
	return nil, service.errLoadingConfiguration
}

func (service *MockService) ApplyConfiguration(interface{}) error {
	return service.errApplyConfiguration
}

func (service *MockService) Restart() error {
	return service.errRestart
}

func (service *MockService) Start() error {
	return service.errStart
}

func (service *MockService) Stop() error {
	return service.errStop
}

var _ = Describe("ServiceStarter", func() {
	It("should start all service", func() {
		reporter := &countEngineReporter{}
		engineStarter := http.NewServiceStarter([]http.Service{
			&MockService{},
		}, reporter)
		err := engineStarter.Start()
		Expect(err).To(BeNil())
		Expect(reporter.countBeforeBegin).To(Equal(1))
		Expect(reporter.countBeforeLoadConfiguration).To(Equal(1))
		Expect(reporter.countAfterLoadConfiguration).To(Equal(1))
		Expect(reporter.countBeforeApplyConfiguration).To(Equal(1))
		Expect(reporter.countAfterApplyConfiguration).To(Equal(1))
		Expect(reporter.countBeforeStart).To(Equal(1))
		Expect(reporter.countAfterStart).To(Equal(1))
	})

	It("should fail loading configuration", func() {
		reporter := &countEngineReporter{}
		engineStarter := http.NewServiceStarter([]http.Service{
			&MockService{
				errLoadingConfiguration: errors.New("loading configuration error"),
			},
		}, reporter)
		err := engineStarter.Start()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("loading configuration error"))
		Expect(reporter.countBeforeBegin).To(Equal(1))
		Expect(reporter.countBeforeLoadConfiguration).To(Equal(1))
		Expect(reporter.countAfterLoadConfiguration).To(Equal(1))
		Expect(reporter.countBeforeApplyConfiguration).To(Equal(0))
		Expect(reporter.countAfterApplyConfiguration).To(Equal(0))
		Expect(reporter.countBeforeStart).To(Equal(0))
		Expect(reporter.countAfterStart).To(Equal(0))
	})

	It("should fail applying configuration", func() {
		reporter := &countEngineReporter{}
		engineStarter := http.NewServiceStarter([]http.Service{
			&MockService{
				errLoadingConfiguration: errors.New("loading configuration error"),
			},
		}, reporter)
		err := engineStarter.Start()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("loading configuration error"))
		Expect(reporter.countBeforeBegin).To(Equal(1))
		Expect(reporter.countBeforeLoadConfiguration).To(Equal(1))
		Expect(reporter.countAfterLoadConfiguration).To(Equal(1))
		Expect(reporter.countBeforeApplyConfiguration).To(Equal(0))
		Expect(reporter.countAfterApplyConfiguration).To(Equal(0))
		Expect(reporter.countBeforeStart).To(Equal(0))
		Expect(reporter.countAfterStart).To(Equal(0))
	})
})
