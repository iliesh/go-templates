### go-templates
Personal Code Templates

#### LOGGER
To import the logger package use:

```golang
go get github.com/iliesh/go-templates/logger/
```

-- Examples:
```golang
log.Info("Info Simple Text of Type/Value: %T/%v", var1, var)
```

-- By default LogLevel is set to Trace, can be changed from main func with:
```golang
log.LogLevel = "info"
```

Valid Levels: Trace/Debug/Info/Warning/Error/Panic

-- By default programname variable will be the name of the executable file, to change it set the ProgramName variable accordingly:
```golang
log.AppName="Program1"
```

-- By default Package Version="0.0.0", to change it - set the following variable:
```golang
log.Version="1.0.0"
```

-- Changing RequestID variable
```golang
log.ReqID = log.RandomString(8)
```
