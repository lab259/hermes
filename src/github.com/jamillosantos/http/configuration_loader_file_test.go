package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"github.com/jamillosantos/http"
	"io/ioutil"
	"fmt"
	"path"
	"time"
)

var _ = Describe("ConfigurationLoaderFile", func() {
	It("should load a file", func() {
		id := fmt.Sprintf("%d", time.Now().Unix())
		dir := os.TempDir()
		fname := path.Join(dir, id)
		ioutil.WriteFile(fname, []byte("this is a file"), 0777)

		loader := http.NewFileConfigurationLoader(dir)
		var data interface{}
		buff, err := loader.Load(id, &data)
		Expect(err).To(BeNil())
		Expect(string(buff)).To(Equal("this is a file"))
	})

	It("should fail unmarshaling a malformed JSON", func() {
		loader := http.NewFileConfigurationLoader("a non existing folder")
		var data interface{}
		buff, err := loader.Load("a non existing file", &data)
		Expect(err).NotTo(BeNil())
		Expect(buff).To(BeNil())
	})
})
