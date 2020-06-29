package log

import (
	"os"
)

// Info log info line
func Info(v ...interface{}) {
	infoLog.Println(v...)
}

// Warning log warning line
func Warning(v ...interface{}) {
	warningLog.Println(v...)
}

// Error log error line
func Error(v ...interface{}) {
	errorLog.Println(v...)
}

// Infof log info line with format
func Infof(format string, v ...interface{}) {
	infoLog.Printf(format, v...)
}

// Warningf log warning line with format
func Warningf(format string, v ...interface{}) {
	warningLog.Printf(format, v...)
}

// Errorf log error line with format
func Errorf(format string, v ...interface{}) {
	errorLog.Printf(format, v...)
}

// Fatal log fatal error and terminate
func Fatal(v ...interface{}) {
	Error(v...)
	Error("Fatal error")
	os.Exit(1)
}

// Close close log file
func Close() {
	Info("Session closed.")

	if file != nil {
		file.Close()
	}
}
