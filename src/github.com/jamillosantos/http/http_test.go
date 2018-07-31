package http_test

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/jamillosantos/macchiato"
	"testing"
	"log"
)

func TestHttp(t *testing.T) {
	log.SetOutput(ginkgo.GinkgoWriter)
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "Http Test Suite")
}
