package config

import (
	"fmt"

	"github.com/mnaufalhilmym/gotracing"
)

func ConfigureTracing(printLevel string, stacktraceLevel string, maxPC uint) {
	printFilter := gotracing.NewLevelFilter(printLevel)
	if printFilter.IsOk() {
		gotracing.SetMinConsolePrintLevel(printFilter.Unwrap())
	} else {
		panic(fmt.Errorf("error reading tracing print filter: invalid value %s", printLevel))
	}

	stacktraceFilter := gotracing.NewLevelFilter(stacktraceLevel)
	if stacktraceFilter.IsOk() {
		gotracing.SetMinStackTrace(stacktraceFilter.Unwrap())
	} else {
		panic(fmt.Errorf("error reading tracing stacktrace filter: invalid value %s", stacktraceLevel))
	}

	gotracing.SetMaxProgramCounters(maxPC)
}
