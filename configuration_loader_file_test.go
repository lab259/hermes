package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"io/ioutil"
	"fmt"
	"path"
	"time"
)

var _ = g.Describe("ConfigurationLoaderFile", func() {
	g.It("should load a file", func() {
		id := fmt.Sprintf("%d", time.Now().Unix())
		dir := os.TempDir()
		fname := path.Join(dir, id)
		ioutil.WriteFile(fname, []byte("this is a file"), 0777)

		loader := NewFileConfigurationLoader(dir)
		var data interface{}
		buff, err := loader.Load(id, &data)
		Expect(err).To(BeNil())
		Expect(string(buff)).To(Equal("this is a file"))
	})

	g.It("should fail unmarshaling a malformed JSON", func() {
		loader := NewFileConfigurationLoader("a non existing folder")
		var data interface{}
		buff, err := loader.Load("a non existing file", &data)
		Expect(err).NotTo(BeNil())
		Expect(buff).To(BeNil())
	})
})
