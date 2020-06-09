package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	codeFile, codeFunc, requestID string
	line                          int
	noLogFile                     bool
)

// Func logRequestID is used to generate a uniq id in order to find the necessary log based on this
func logRequestID() (id string) {

	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("Logger: Unable to generate a request id: %v\n", err)
		return "00000000"
	}
	return fmt.Sprintf("%x", b)
}

// Func runtimeInfo is used to get runtime information like file name, function and line number of executed script
func runtimeInfo(depthList ...int) (cf, fct string, l int) {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, ok := runtime.Caller(depth)

	if !ok {
		fmt.Printf("Logger: Unable to get runtime data\n")
		return "?", "runtimeInfo", 0
	}

	return path.Base(file), runtime.FuncForPC(function).Name(), line
}

func logStdOut(level, msg string) {
	logFormat := logrus.New()
	logFormat.SetLevel(logrus.TraceLevel)
	logFormat.SetOutput(os.Stdout)

	logFormat.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: false,
		DisableSorting:         false,
		DisableColors:          false,
		ForceColors:            true,
		PadLevelText:           false,
	})

	logPrint(logFormat, level, msg)

}

func logFile(level, msg string) {
	logFormat := logrus.New()
	logFormat.SetLevel(logrus.TraceLevel)
	logFormat.SetOutput(os.Stdout)

	logPath := "/var/log/scripts"

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.Mkdir(logPath, 0666); err != nil {
			logFormat.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:          true,
				TimestampFormat:        "2006-01-02 15:04:05",
				DisableLevelTruncation: false,
			})
			logFormat.WithFields(logrus.Fields{
				"version":   version,
				"requestid": requestID,
				"function":  codeFunc,
				"file":      codeFile,
				"line":      line,
			}).Warn("Failed to create Log Folder, Error: " + fmt.Sprintf("%v", err))
			noLogFile = true
			return
		}
	}

	logfile, err := os.OpenFile(logPath+"/"+programName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		logFormat.SetOutput(logfile)
	} else {
		codeFile, codeFunc, line := runtimeInfo(1)
		logFormat.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			DisableLevelTruncation: true,
		})
		logFormat.WithFields(logrus.Fields{
			"version":   version,
			"requestid": requestID,
			"function":  codeFunc,
			"file":      codeFile,
			"line":      line,
		}).Warn("Failed to open log file, Error: " + fmt.Sprintf("%v", err))
		noLogFile = true
		return
	}
	defer logfile.Close()

	logFormat.SetNoLock()
	logFormat.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		DisableSorting:         false,
		DisableColors:          true,
		ForceQuote:             true,
	})

	logPrint(logFormat, level, msg)
}

func logPrint(logFormat *logrus.Logger, level, msg string) {

	log := logFormat.WithFields(logrus.Fields{
		"version":   version,
		"file":      codeFile,
		"requestid": requestID,
		"function":  codeFunc,
		"line":      line,
	})

	switch level {
	case "trace":
		log.Trace(msg)
	case "debug":
		log.Debug(msg)
	case "info":
		log.Info(msg)
	case "warning":
		log.Warning(msg)
	case "error":
		log.Error(msg)
	case "fatal":
		log.Fatal(msg)
	case "panic":
		log.Panic(msg)
	default:
		log.Trace(msg)
	}
}

func logger(level, msg string) {

	if requestID == "" {
		requestID = logRequestID()
	}

	codeFile, codeFunc, line = runtimeInfo(2)
	logStdOut(strings.ToLower(level), msg)

	if !noLogFile {
		logFile(strings.ToLower(level), msg)
	}

}
