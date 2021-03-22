package logger

// Version 2.0.1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	// AppName - Name of the running application
	AppName = ""
	// Version - Application Version
	Version = "0.0.0"
	// LogLevel - Application loglevel
	LogLevel = "trace"
	// ReqID - Unique ID
	ReqID = ""
	// Color - Colorize log output
	Color = true
	// Env variable - Setting the running environment (prod/dev)
	// For the prod environment - the log will be formatted automatically as a json object and colors
	// will be disabled
	Env = "dev"

	levelColor = map[string]string{
		"TRACE": colorWhite,
		"DEBUG": colorGreen,
		"INFO ": colorCyan,
		"WARN ": colorYellow,
		"ERROR": colorRed,
		"PANIC": colorRedBg,
	}

	logLevelNumber = map[string]int{
		"trace":   6,
		"debug":   5,
		"info":    4,
		"warning": 3,
		"error":   2,
		"panic":   1,
	}
)

type logT struct {
	Time      string `json:"time,omitempty"`
	Level     string `json:"level,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Msg       string `json:"msg,omitempty"`
	File      string `json:"file,omitempty"`
	Func      string `json:"func,omitempty"`
	Line      string `json:"line,omitempty"`
	AppName   string `json:"app_name,omitempty"`
	Version   string `json:"version,omitempty"`
}

const (
	colorReset = "\033[0m"
	// colorBlack  = "\033[30m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	// colorBlue   = "\033[34m"
	// colorPurple = "\033[35m"
	colorCyan  = "\033[36m"
	colorWhite = "\033[37m"
	colorGray  = "\033[90m"

	colorRedBg = "\033[41m"

	charBytes     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ReqID = RandomString(8)
	AppName = path.Base(os.Args[0])
}

func RandomString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(charBytes) {
			b[i] = charBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Func runtimeInfo is used to get runtime information like file name, function and line number of executed script
func runtimeInfo(depthList ...int) (cf, fct, l string) {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, ok := runtime.Caller(depth)

	if !ok {
		fmt.Printf("Logger: Unable to get runtime data\n")
		return "?", "runtimeInfo", "0"
	}

	return path.Base(file), runtime.FuncForPC(function).Name(), strconv.Itoa(line)
}

func Trace(format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 6 {
		return
	}

	output := logT{}
	output.Level = "TRACE"
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func Debug(format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 5 {
		return
	}

	output := logT{}
	output.Level = "DEBUG"
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func Info(format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 4 {
		return
	}

	output := logT{}
	output.Level = "INFO "
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func Warning(format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 3 {
		return
	}

	output := logT{}
	output.Level = "WARN "
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func Error(format string, s ...interface{}) {

	if logLevelNumber[strings.ToLower(LogLevel)] < 2 {
		return
	}

	output := logT{}
	output.Level = "ERROR"
	output.Msg = format

	o, err := logFormat(output)
	if err != nil {
		fmt.Printf("[LOGGER] - Error Formatting the Object: %s\n", err.Error())
		return
	}

	fmt.Printf(o, s...)
}

func Panic(format string, s ...interface{}) {
	output := logT{}
	output.Level = "PANIC"
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

func logFormat(logData logT) (string, error) {
	codeFile, codeFunc, line := runtimeInfo(3)
	currentTime := time.Now()
	timeFmt := currentTime.Format("2006.01.02 15:04:05")

	if Env == "prod" {
		logData.Time = timeFmt
		logData.RequestID = ReqID
		logData.File = codeFile
		logData.Func = codeFunc
		logData.Line = line
		logData.AppName = AppName
		logData.Version = Version
		out, err := json.Marshal(logData)
		if err != nil {
			fmt.Printf("[LOGGER] - Unable to Encode to JSON\n")
			return "", err
		}
		return string(out) + "\n", nil
	}

	if Color {
		logData.Time = levelColor[logData.Level] + "[" + timeFmt + "] " + colorReset
		logData.RequestID = "[" + ReqID + "]  "
		logData.Msg = levelColor[logData.Level] + logData.Msg + colorReset + "    "
		logData.File = colorGray + codeFile + "," + colorReset
		logData.Func = colorGray + codeFunc + "," + colorReset
		logData.Line = colorGray + "(" + line + ")" + "," + colorReset
		logData.AppName = colorGray + "appname:" + AppName + "," + colorReset
		logData.Version = colorGray + "version:" + Version + colorReset
		logData.Level = levelColor[logData.Level] + "[" + logData.Level + "] " + colorReset
	} else {
		logData.Time = "[" + timeFmt + "] "
		logData.RequestID = "[" + ReqID + "]  "
		logData.Msg = logData.Msg + "    "
		logData.File = codeFile + ","
		logData.Func = codeFunc + ","
		logData.Line = "(" + line + ")" + ","
		logData.AppName = "appname: " + AppName + ","
		logData.Version = "version:" + Version
		logData.Level = "[" + logData.Level + "] "
	}
	out := logData.Time + logData.Level + logData.RequestID + logData.Msg + logData.File + logData.Func + logData.Line + logData.AppName + logData.Version + "\n"

	return out, nil
}
