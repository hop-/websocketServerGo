package log

import (
	"io"
	"log"
	"os"
)

var (
	file       *os.File
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
)

// Init initiate output streams
func Init() {
	infoLog = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime)
	warningLog = log.New(os.Stderr, "Warning: ", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime)
}

// InitFile initiate output streams and setup log file
func InitFile(fileName string) {
	var err error

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	stdOutWriter := io.MultiWriter(file, os.Stdout)
	stdErrWriter := io.MultiWriter(file, os.Stderr)
	infoLog = log.New(stdOutWriter, "Info: ", log.Ldate|log.Ltime)
	warningLog = log.New(stdErrWriter, "Warning: ", log.Ldate|log.Ltime)
	errorLog = log.New(stdErrWriter, "Error: ", log.Ldate|log.Ltime)
}
