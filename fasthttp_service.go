package hermes

import (
	"errors"

	rscsrv "github.com/lab259/go-rscsrv"
	"github.com/valyala/fasthttp"
)

// FasthttpServiceConfiguration keeps all the configuration needed to start the
// `FasthttpService`.
type FasthttpServiceConfiguration struct {
	Bind string
	TLS  *FasthttpServiceConfigurationTLS
}

// FasthttpServiceConfigurationTLS keeps the configuration for starting a TLS
// server.
type FasthttpServiceConfigurationTLS struct {
	CertFile string
	KeyFile  string
}

// FasthttpService implements the server for starting
type FasthttpService struct {
	running       bool
	Configuration FasthttpServiceConfiguration
	Server        fasthttp.Server
}

// LoadConfiguration does not do anything in this implementation. This methods
// is just a placeholder to be overwritten on its usage.
func (service *FasthttpService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ApplyConfiguration checks if the passing interface is a
// `FasthttpServiceConfiguration` and applies its configuration to the service.
func (service *FasthttpService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case *FasthttpServiceConfiguration:
		service.Configuration = *c
		return nil
	case FasthttpServiceConfiguration:
		service.Configuration = c
		return nil
	}
	return rscsrv.ErrWrongConfigurationInformed
}

// Restart returns an error due to fasthttp not being able to stop the service.
func (service *FasthttpService) Restart() error {
	if service.running {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start ListenAndServe the server. This method is blocking because it uses
// the fasthttp.ListenAndServe implementation.
func (service *FasthttpService) Start() error {
	service.running = true

	if service.Configuration.TLS == nil {
		return service.Server.ListenAndServe(service.Configuration.Bind)
	}

	return service.Server.ListenAndServeTLS(service.Configuration.Bind, service.Configuration.TLS.CertFile, service.Configuration.TLS.KeyFile)
}

// Stop closes the listener and waits the `Start` to stop.
func (service *FasthttpService) Stop() error {
	if service.running {
		service.running = false
		return service.Server.Shutdown()
	}
	return nil
}
