package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jamillosantos/http"
)

type UnmarshalingTest struct {
	Name1 string
	Name2 int
}

var _ = Describe("ConfigurationUnmarshelerJson", func() {
	It("should unmarshal a configuration", func() {
		var dst UnmarshalingTest
		Expect(http.DefaultConfigurationUnmarshalerJson.Unmarshal([]byte(`{"name1": "value 1", "name2": 2}`), &dst)).To(BeNil())
		Expect(dst.Name1).To(Equal("value 1"))
		Expect(dst.Name2).To(Equal(2))
	})

	It("should fail unmarshaling a malformed JSON", func() {
		var dst UnmarshalingTest
		Expect(http.DefaultConfigurationUnmarshalerJson.Unmarshal([]byte(`this is not a JSON`), &dst)).NotTo(BeNil())
	})
})
