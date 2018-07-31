package http

// ConfigurationUnmarshaler describes the unmarshaling contract of a configuration.
type ConfigurationUnmarshaler interface {
	Unmarshal(buf []byte, dst interface{}) error
}
