package http

// ConfigurationLoader defines the contract to load a configuration
// from a repository.
//
// The repository is an abstract idea that can be represented as a
// directory, a S3 bucket, or "anything" else.
type ConfigurationLoader interface {
	// Load receives the `id` of the configuration and Unmarshals it
	// int the `dst` pointer. If no error is reported the method will
	// return nil otherwise the error will be returned.
	Load(id string) ([]byte, error)
}
