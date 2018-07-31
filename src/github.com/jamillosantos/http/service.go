package http

// Service is an abstraction for implementing parts that can be loaded,
// reloaded, started and stopped inside of the system.
//
// Maybe you can implement your HTTP service like this, or your Redis resource.
// As simple and wide as it could be this directive will provide an defined
// signature to implement all your resources.
type Service interface {
	Load() error
	Reload() error
	Start() error
	Stop() error
}
