package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jamillosantos/http"
)

type UnmarshalingYamlTest struct {
	Name1 string
	Name2 int
}

var _ = Describe("ConfigurationUnmarshelerYaml", func() {
	It("should unmarshal a yaml", func() {
		var dst UnmarshalingYamlTest
		Expect(http.DefaultConfigurationUnmarshalerYaml.Unmarshal([]byte(`name1: "value 1"
name2: 2
`), &dst)).To(BeNil())
		Expect(dst.Name1).To(Equal("value 1"))
		Expect(dst.Name2).To(Equal(2))
	})

	It("should fail unmarshaling a malformed YAML", func() {
		var dst UnmarshalingTest
		Expect(http.DefaultConfigurationUnmarshalerYaml.Unmarshal([]byte(`this is not an YAML`), &dst)).NotTo(BeNil())
	})
})
