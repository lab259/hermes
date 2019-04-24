package http

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"
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
		buff, err := loader.Load(id)
		Expect(err).To(BeNil())
		Expect(string(buff)).To(Equal("this is a file"))
	})

	g.It("should fail unmarshaling a malformed JSON", func() {
		loader := NewFileConfigurationLoader("a non existing folder")
		buff, err := loader.Load("a non existing file")
		Expect(err).NotTo(BeNil())
		Expect(buff).To(BeNil())
	})
})
