# go-templates
Personal Code Templates

### LOGGER ###
# To import the logger package use:
go get github.com/iliesh/go-templates/logger/

-- Examples:
log.Info("Info Simple Text of Type/Value: %T/%v", var1, var)

Valid Levels: Trace/Debug/Info/Warning/Error/Fatal/Panic

-- To Disable File Logging, set NoLogFile variable to true 
logger.NoLogFile=true

-- By default programname variable will be the name of the executable file, to change it set the ProgramName variable accordingly:
logger.ProgramName="Program1"

-- By default Package Version="0.0.0", to change it - set the following variable:
logger.Version="1.0.0"
