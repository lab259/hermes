package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type UnmarshalingYamlTest struct {
	Name1 string
	Name2 int
}

var _ = g.Describe("ConfigurationUnmarshalerYaml", func() {
	g.It("should unmarshal a yaml", func() {
		var dst UnmarshalingYamlTest
		Expect(DefaultConfigurationUnmarshalerYaml.Unmarshal([]byte(`name1: "value 1"
name2: 2
`), &dst)).To(BeNil())
		Expect(dst.Name1).To(Equal("value 1"))
		Expect(dst.Name2).To(Equal(2))
	})

	g.It("should fail unmarshaling a malformed YAML", func() {
		var dst UnmarshalingTest
		Expect(DefaultConfigurationUnmarshalerYaml.Unmarshal([]byte(`this is not an YAML`), &dst)).NotTo(BeNil())
	})
})
