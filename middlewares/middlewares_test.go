package middlewares_test

import (
	"log"
	"os"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

func TestMiddlewares(t *testing.T) {
	log.SetOutput(ginkgo.GinkgoWriter)
	gomega.RegisterFailHandler(ginkgo.Fail)

	description := "Http/Middleware Test Suite"
	if os.Getenv("CI") == "" {
		macchiato.RunSpecs(t, description)
	} else {
		reporterOutputDir := "../test-results"
		os.MkdirAll(reporterOutputDir, os.ModePerm)
		junitReporter := reporters.NewJUnitReporter("../test-results/middlewares.xml")
		macchiatoReporter := macchiato.NewReporter()
		ginkgo.RunSpecsWithCustomReporters(t, description, []ginkgo.Reporter{macchiatoReporter, junitReporter})
	}
}
