package middlewares_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/rlog"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

func TestMiddlewares(t *testing.T) {
	rlog.SetOutput(ginkgo.GinkgoWriter)
	log.SetOutput(ginkgo.GinkgoWriter)
	gomega.RegisterFailHandler(ginkgo.Fail)

	description := "Hermes/Middleware Test Suite"
	if os.Getenv("CI") == "" {
		macchiato.RunSpecs(t, description)
	} else {
		reporterOutputDir := "../test-results/hermes/middlewares"
		os.MkdirAll(reporterOutputDir, os.ModePerm)
		junitReporter := reporters.NewJUnitReporter(path.Join(reporterOutputDir, "results.xml"))
		macchiatoReporter := macchiato.NewReporter()
		ginkgo.RunSpecsWithCustomReporters(t, description, []ginkgo.Reporter{macchiatoReporter, junitReporter})
	}
}
