package http

import "errors"

var (
	ErrWrongConfigurationInformed = errors.New("wrong configuration informed")
	ErrServiceNotRunning          = errors.New("service not running")
)

// Service is an abstraction for implementing parts that can be loaded,
// reloaded, started and stopped inside of the system.
//
// Maybe you can implement your HTTP service like this, or your Redis resource.
// As simple and wide as it could be this directive will provide an defined
// signature to implement all your resources.
type Service interface {
	// Name identifies the service.
	Name() string

	// Loads the configuration. If successful nil will be returned, otherwise
	// the error.
	LoadConfiguration() (interface{}, error)

	// Applies a given configuration object to the service. If successful nil
	// will be returned, otherwise the error.
	ApplyConfiguration(interface{}) error

	// Restarts the service. If successful nil will be returned, otherwise the
	// error.
	Restart() error

	// Start starts the service. If successful nil will be returned, otherwise
	// the error.
	Start() error

	// Stop stops the service. If successful nil will be returned, otherwise the
	// error.
	Stop() error
}
