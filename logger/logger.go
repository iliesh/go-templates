package logger

// Package Version 1.0.1

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// ProgramName Variable
	ProgramName string
	// Version Variable
	Version string = "0.0.0"
	// NoLogFile Variable
	NoLogFile          bool
	codeFile, codeFunc string
	line               int
	requestID          int64
)

func init() {
	// Generate RequestID value
	if requestID == 0 {
		requestID = time.Now().UnixNano()
	}
	ProgramName = path.Base(os.Args[0])
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

func logStdOut(level, format string, args ...interface{}) {
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

	logPrint(logFormat, level, format, args...)

}

func logFile(level, format string, args ...interface{}) {
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
				"version":   Version,
				"requestid": requestID,
				"function":  codeFunc,
				"file":      codeFile,
				"line":      line,
			}).Warnf("Failed to create Log Folder, Error: %v", err)
			NoLogFile = true
			return
		}
	}

	logfile, err := os.OpenFile(logPath+"/"+ProgramName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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
			"version":   Version,
			"requestid": requestID,
			"function":  codeFunc,
			"file":      codeFile,
			"line":      line,
		}).Warnf("Failed to open log file, Error: %v", err)
		NoLogFile = true
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

	logPrint(logFormat, level, format, args...)
}

func logPrint(logFormat *logrus.Logger, level, format string, args ...interface{}) {

	codeFile, codeFunc, line = runtimeInfo(4)

	log := logFormat.WithFields(logrus.Fields{
		"version":   Version,
		"file":      codeFile,
		"requestid": requestID,
		"function":  codeFunc,
		"line":      line,
	})

	switch level {
	case "trace":
		log.Tracef(format, args...)
	case "debug":
		log.Debugf(format, args...)
	case "info":
		log.Infof(format, args...)
	case "warning":
		log.Warnf(format, args...)
	case "error":
		log.Errorf(format, args...)
	case "fatal":
		log.Fatalf(format, args...)
	case "panic":
		log.Panicf(format, args...)
	default:
		log.Tracef(format, args...)
	}
}

// Trace logs a message at level Trace on the standard logger and to the file.
func Trace(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("trace", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("trace", format, args...)
	}
}

// Debug logs a message at level Debug on the standard logger and to the file.
func Debug(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("debug", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("debug", format, args...)
	}
}

// Info logs a message at level Info on the standard logger and to the file.
func Info(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("info", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("info", format, args...)
	}
}

// Warn logs a message at level Warn on the standard logger and to the file.
func Warn(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("warning", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("warning", format, args...)
	}
}

// Error logs a message at level Error on the standard logger and to the file.
func Error(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("error", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("error", format, args...)
	}
}

// Panic logs a message at level Panic on the standard logger and to the file.
func Panic(format string, args ...interface{}) {
	// Logging to stdout
	logStdOut("panic", format, args...)

	if !NoLogFile {
		// Logging to file
		logFile("panic", format, args...)
	}
}
