package http

import "gopkg.in/yaml.v2"

type ConfigurationUnmarshalerYaml struct {
}

var DefaultConfigurationUnmarshalerYaml ConfigurationUnmarshalerYaml

// Unmarshal is an abstract method that should be override
func (loader *ConfigurationUnmarshalerYaml) Unmarshal(buff []byte, dst interface{}) error {
	return yaml.Unmarshal(buff, dst)
}
