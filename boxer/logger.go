package main

import (
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func getLogLevel() pterm.LogLevel {
	if traceEnabled() {
		return pterm.LogLevelTrace
	}
	if debugEnabled() {
		return pterm.LogLevelDebug
	}
	return pterm.LogLevelInfo
}

var logger = pterm.DefaultLogger.WithLevel(getLogLevel())

func debugEnabled() bool {
	debugEnvVar := os.Getenv("UVBOX_DEBUG")
	debugEnvVar = strings.ToLower(debugEnvVar)
	return debugEnvVar == "1" || debugEnvVar == "true"
}

func traceEnabled() bool {
	debugEnvVar := os.Getenv("UVBOX_TRACE")
	debugEnvVar = strings.ToLower(debugEnvVar)
	return debugEnvVar == "1" || debugEnvVar == "true"
}
