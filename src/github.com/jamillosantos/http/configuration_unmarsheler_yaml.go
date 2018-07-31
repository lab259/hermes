package http

import "gopkg.in/yaml.v2"

type ConfigurationUnmarshelerYaml struct {
}

// Unmarshal is an abstract method that should be override
func (loader *ConfigurationUnmarshelerYaml) Unmarshal(buff []byte, dst interface{}) error {
	return yaml.Unmarshal(buff, dst)
}
