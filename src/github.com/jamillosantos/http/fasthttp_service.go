package http

import (
	"github.com/valyala/fasthttp"
	"errors"
	"net"
	"sync"
	"io"
)

// FasthttpServiceConfiguration keeps all the configuration needed to start the
// `FasthttpService`.
type FasthttpServiceConfiguration struct {
	Bind string
}

// FasthttpService implements the server for starting
type FasthttpService struct {
	running       bool
	waitToFinish  sync.WaitGroup
	Router        *Router
	Configuration FasthttpServiceConfiguration
	Listener      net.Listener
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
	return ErrWrongConfigurationInformed
}

// Reload returns an error due to fasthttp not being able to stop the service.
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
	service.Server.Handler = service.Router.Handler
	ln, err := net.Listen("tcp4", service.Configuration.Bind)
	if err != nil {
		return err
	}
	service.Listener = ln
	service.running = true
	err = service.Server.Serve(ln)
	defer service.waitToFinish.Done()
	if err == io.EOF && !service.running {
		return nil
	}
	return err
}

// Stop closes the listener and waits the `Start` to stop.
func (service *FasthttpService) Stop() error {
	if service.running {
		service.waitToFinish.Add(1)
		err := service.Listener.Close()
		if err != nil {
			return err
		}
		service.running = false
		service.waitToFinish.Wait()
	}
	return nil
}
