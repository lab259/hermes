package hermes

import (
	"sync"
	"time"

	"github.com/lab259/go-rscsrv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hermes", func() {
	Describe("Application", func() {
		It("should start and stop a app", func(done Done) {
			app := NewApplication(ApplicationConfig{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(RouterConfig{}))
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()

				wg.Done()
				Expect(app.Start()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 500)
			wg.Add(1)
			Expect(app.Stop()).To(BeNil())
			wg.Wait()
			done <- true
		}, 1)

		It("should start and stop a app with services", func(done Done) {
			var serviceA testService
			var serviceB testService

			serviceStarter := rscsrv.QuietServiceStarter(&serviceA, &serviceB)
			Expect(serviceStarter.Start()).To(Succeed())

			Expect(serviceA.running).To(BeTrue())
			Expect(serviceB.running).To(BeTrue())

			app := NewApplication(ApplicationConfig{
				Name:           "Testing",
				ServiceStarter: serviceStarter,
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(RouterConfig{}))
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()

				wg.Done()
				Expect(app.Start()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 500)
			wg.Add(1)
			Expect(app.Stop()).To(BeNil())
			wg.Wait()

			Expect(serviceA.running).To(BeFalse())
			Expect(serviceB.running).To(BeFalse())
			done <- true
		}, 1)

		It("should restart a app", func(done Done) {
			app := NewApplication(ApplicationConfig{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, DefaultRouter())

			ch := make(chan string, 10)

			// Just ignore result
			go func() {
				ch <- "service:step1:begin"
				app.Start()
				ch <- "service:step1:end"
			}()

			time.Sleep(time.Millisecond * 500) // Waits for the service a bit

			go func() {
				ch <- "service:step2:begin"
				app.Restart()
				time.Sleep(time.Millisecond * 500)
				ch <- "service:step2:end"
			}()

			time.Sleep(time.Millisecond * 500) // Waits for the service a bit

			Expect(app.Stop()).To(BeNil())
			Expect(<-ch).To(Equal("service:step1:begin"))
			Expect(<-ch).To(Equal("service:step2:begin"))
			Expect(<-ch).To(Equal("service:step1:end"))
			Expect(<-ch).To(Equal("service:step2:end"))
			done <- true
		}, 2)

		It("should return name", func() {
			app := NewApplication(ApplicationConfig{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(RouterConfig{}))
			Expect(app.Name()).To(Equal("Testing"))

			app2 := NewApplication(ApplicationConfig{
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(RouterConfig{}))
			Expect(app2.Name()).To(Equal("Application"))
		})

		It("should fail with misconfiguration", func() {
			app := NewApplication(ApplicationConfig{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":FAIL",
				},
			}, NewRouter(RouterConfig{}))
			Expect(app.Start()).ToNot(BeNil())
		})
	})
})

type testService struct {
	running bool
}

func (service *testService) Name() string {
	return "Service A"
}

func (service *testService) LoadConfiguration() (interface{}, error) {
	return nil, nil
}

func (service *testService) ApplyConfiguration(interface{}) error {
	return nil
}

func (service *testService) Restart() error {
	if err := service.Stop(); err != nil {
		return err
	}
	return service.Start()
}

func (service *testService) Start() error {
	service.running = true
	return nil
}

func (service *testService) Stop() error {
	service.running = false
	return nil
}
