package http

import "encoding/json"

type ConfigurationUnmarshelerJson struct {
}

var DefaultConfigurationUnmarshalerJson ConfigurationUnmarshelerJson

// Unmarshal is an abstract method that should be override
func (loader *ConfigurationUnmarshelerJson) Unmarshal(buff []byte, dst interface{}) error {
	return json.Unmarshal(buff, dst)
}
