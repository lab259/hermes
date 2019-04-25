package http

import (
	"sync"
	"time"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = describe("Http", func() {
	describe("Application", func() {
		it("should start and stop a app", func(done g.Done) {
			app := NewApplication(Config{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(nil))
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer g.GinkgoRecover()

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

		it("should restart a app", func(done g.Done) {
			app := NewApplication(Config{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(nil))

			ch := make(chan string, 10)

			// Just ignore result
			go func() {
				ch <- "service:step1:begin"
				app.Start()
				ch <- "service:step1:end"
			}()

			time.Sleep(time.Millisecond * 50) // Waits for the service a bit

			go func() {
				ch <- "service:step2:begin"
				app.Restart()
				time.Sleep(time.Millisecond * 50)
				ch <- "service:step2:end"
			}()

			time.Sleep(time.Millisecond * 50) // Waits for the service a bit

			Expect(app.Stop()).To(BeNil())
			Expect(<-ch).To(Equal("service:step1:begin"))
			Expect(<-ch).To(Equal("service:step2:begin"))
			Expect(<-ch).To(Equal("service:step1:end"))
			Expect(<-ch).To(Equal("service:step2:end"))
			done <- true
		}, 0.5)

		it("should return name", func() {
			app := NewApplication(Config{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(nil))
			Expect(app.Name()).To(Equal("Testing"))

			app2 := NewApplication(Config{
				HTTP: FasthttpServiceConfiguration{
					Bind: ":0",
				},
			}, NewRouter(nil))
			Expect(app2.Name()).To(Equal("Application"))
		})

		it("should fail with misconfiguration", func() {
			app := NewApplication(Config{
				Name: "Testing",
				HTTP: FasthttpServiceConfiguration{
					Bind: ":FAIL",
				},
			}, NewRouter(nil))
			Expect(app.Start()).ToNot(BeNil())
		})
	})
})
