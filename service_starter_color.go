package http

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	colorSuccess   = color.New(color.FgGreen)
	colorError     = color.New(color.FgRed)
	colorWarning   = color.New(color.FgYellow)
	colorHighlight = color.New(color.FgWhite, color.Bold)
	colorMuted     = color.New(color.FgHiBlack)

	colorBold = color.New(color.Bold)
)

// Success formats text as
var (
	formatSuccess   = colorSuccess.SprintfFunc()
	formatError     = colorError.SprintfFunc()
	formatWarning   = colorWarning.SprintfFunc()
	formatHighlight = colorHighlight.SprintfFunc()
	formatMuted     = colorMuted.SprintfFunc()

	formatBold = colorSuccess.SprintfFunc()
)

type ColorServiceReporter struct{}

func (*ColorServiceReporter) BeforeBegin(service Service) {
	fmt.Printf("%s\n", formatHighlight(service.Name()))
}

const colorTitleL1 string = "    %-27s"

func (*ColorServiceReporter) BeforeLoadConfiguration(service Service) {
	fmt.Printf(colorTitleL1, "Loading configuration ...")
}

func printError(err error) {
	var t string
	if err != nil {
		t = formatBold(formatSuccess("Error"))
		fmt.Printf("[%s]\n    %s\n", t, err.Error())
		return
	}
	t = formatBold(formatSuccess("OK"))
	fmt.Printf("[%s]\n", t)
}

func (*ColorServiceReporter) AfterLoadConfiguration(service Service, conf interface{}, err error) {
	printError(err)
}

func (*ColorServiceReporter) BeforeApplyConfiguration(service Service) {
	fmt.Printf(colorTitleL1, "Applying configuration ...")
}

func (*ColorServiceReporter) AfterApplyConfiguration(service Service, conf interface{}, err error) {
	printError(err)
}

func (*ColorServiceReporter) BeforeStart(service Service) {
	fmt.Printf(colorTitleL1, "Starting ...")
}

func (*ColorServiceReporter) AfterStart(service Service, err error) {
	printError(err)
}
