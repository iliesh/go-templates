package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func TraceX(id, format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 6 {
		return
	}

	output := logT{}
	output.Level = "TRACE"
	output.RequestID = id
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func DebugX(id, format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 5 {
		return
	}

	output := logT{}
	output.Level = "DEBUG"

	if Env == "prod" {
		output.RequestID = id
	} else {
		output.RequestID = "[" + id + "]  "
	}

	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func InfoX(id, format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 4 {
		return
	}

	output := logT{}
	output.Level = "INFO "
	output.RequestID = id
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func WarningX(id, format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 3 {
		return
	}

	output := logT{}
	output.Level = "WARN "
	output.RequestID = id
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func ErrorX(id, format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 2 {
		return
	}

	output := logT{}
	output.Level = "ERROR"
	output.RequestID = id
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func PanicX(id, format string, s ...interface{}) {
	output := logT{}
	output.Level = "PANIC"
	output.RequestID = id
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)

	// Print Stack Trace
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	fmt.Printf("%s", buf)
	os.Exit(1)
}
