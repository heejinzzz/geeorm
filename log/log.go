package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mutex    sync.Mutex
)

// log methods
var (
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Error  = errorLog.Println
	Errorf = errorLog.Printf
)

type Level int

// log levels
const (
	InfoLevel  Level = 0
	ErrorLevel Level = 1
	Disabled   Level = 2
)

// SetLevel controls log level
func SetLevel(level Level) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
}
