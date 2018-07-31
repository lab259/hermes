package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = g.Describe("Services", func() {
	g.Describe("Fasthttp Service", func() {
		g.It("should not return any error loading the service", func() {
			var service FasthttpService
			result, err := service.LoadConfiguration()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("not implemented"))
			Expect(result).To(BeNil())
		})

		g.It("should apply a given pointer configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(&FasthttpServiceConfiguration{
				Bind: "12345",
			})).To(BeNil())
			Expect(service.Configuration.Bind).To(Equal("12345"))
		})

		g.It("should apply a given configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(FasthttpServiceConfiguration{
				Bind: "12345",
			})).To(BeNil())
			Expect(service.Configuration.Bind).To(Equal("12345"))
		})

		g.It("should fail applying a wrong type configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(map[string]interface{}{
				"bind": "12345",
			})).To(Equal(ErrWrongConfigurationInformed))
		})

		g.It("should start and stop the service", func(done g.Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer g.GinkgoRecover()

				wg.Done()
				Expect(service.Start()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 50)
			wg.Add(1)
			Expect(service.Stop()).To(BeNil())
			wg.Wait()
			done <- true
		}, 0.5)

		g.It("should restart the service that is not started", func(done g.Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer g.GinkgoRecover()

				wg.Done()
				Expect(service.Restart()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 50)
			wg.Add(1)
			Expect(service.Stop()).To(BeNil())
			wg.Wait()
			done <- true
		}, 0.5)

		g.It("should stop a stopped service", func() {
			var service FasthttpService
			Expect(service.Stop()).To(BeNil())
		})

		g.It("should restart the service", func(done g.Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port

			ch := make(chan string, 10)

			// Just ignore result
			go func() {
				ch <- "service:step1:begin"
				service.Start()
				ch <- "service:step1:end"
			}()

			time.Sleep(time.Millisecond * 50) // Waits for the service a bit

			go func() {
				ch <- "service:step2:begin"
				service.Restart()
				ch <- "service:step2:end"
			}()

			time.Sleep(time.Millisecond * 50) // Waits for the service a bit

			Expect(service.Stop()).To(BeNil())
			Expect(<-ch).To(Equal("service:step1:begin"))
			Expect(<-ch).To(Equal("service:step2:begin"))
			Expect(<-ch).To(Equal("service:step1:end"))
			Expect(<-ch).To(Equal("service:step2:end"))
			done <- true
		}, 0.5)
	})
})
