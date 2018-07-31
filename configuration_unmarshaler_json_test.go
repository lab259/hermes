package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type UnmarshalingTest struct {
	Name1 string
	Name2 int
}

var _ = g.Describe("ConfigurationUnmarshelerJson", func() {
	g.It("should unmarshal a configuration", func() {
		var dst UnmarshalingTest
		Expect(DefaultConfigurationUnmarshalerJson.Unmarshal([]byte(`{"name1": "value 1", "name2": 2}`), &dst)).To(BeNil())
		Expect(dst.Name1).To(Equal("value 1"))
		Expect(dst.Name2).To(Equal(2))
	})

	g.It("should fail unmarshaling a malformed JSON", func() {
		var dst UnmarshalingTest
		Expect(DefaultConfigurationUnmarshalerJson.Unmarshal([]byte(`this is not a JSON`), &dst)).NotTo(BeNil())
	})
})
