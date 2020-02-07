package hermes

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"time"
)

var _ = Describe("Hermes", func() {
	Describe("Service", func() {
		It("should not start service in a used", func() {

			serverPort := ":2699"

			firstRouter := DefaultRouter()
			firstRouter.Get("/", func(_ Request, res Response) Result {
				return res.Data(map[string]string{"test": "true"})
			})

			firstServer := NewApplication(ApplicationConfig{
				Name: "Running test app",
				HTTP: FasthttpServiceConfiguration{
					Bind: serverPort,
				},
			}, firstRouter)

			go func() {
				defer GinkgoRecover()
				Expect(firstServer.Start()).ShouldNot(HaveOccurred())
			}()

			time.Sleep(10 * time.Millisecond)

			router := DefaultRouter()
			router.Get("/", func(_ Request, res Response) Result {
				return res.Data(map[string]string{"test": "false"})
			})

			server := NewApplication(ApplicationConfig{
				Name: "Should brake test app",
				HTTP: FasthttpServiceConfiguration{
					Bind: serverPort,
				},
			}, router)

			err := server.Start()

			var opErr *net.OpError

			Expect(err).To(BeAssignableToTypeOf(opErr))
			opErr = err.(*net.OpError)

			Expect(opErr.Op).To(Equal("listen"))

		})
	})
})
