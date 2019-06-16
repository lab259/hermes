package hermes

import (
	"log"
	"os"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/lab259/rlog"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	rlog.SetOutput(ginkgo.GinkgoWriter)
	log.SetOutput(ginkgo.GinkgoWriter)
	gomega.RegisterFailHandler(ginkgo.Fail)

	description := "Hermes Test Suite"
	if os.Getenv("CI") == "" {
		macchiato.RunSpecs(t, description)
	} else {
		reporterOutputDir := "./test-results"
		os.MkdirAll(reporterOutputDir, os.ModePerm)
		junitReporter := reporters.NewJUnitReporter("./test-results/hermes.xml")
		macchiatoReporter := macchiato.NewReporter()
		ginkgo.RunSpecsWithCustomReporters(t, description, []ginkgo.Reporter{macchiatoReporter, junitReporter})
	}
}
