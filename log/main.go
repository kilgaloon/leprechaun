package log

import (
	"fmt"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

// Logs struct holds path to different logs
type Logs struct {
	Debug    bool
	ErrorLog string
	InfoLog  string
}

// Error logs everything bad that happens in application
func (l Logs) Error(message string, v ...interface{}) {
	if len(l.ErrorLog) > 0 {
		file, err := os.OpenFile(l.ErrorLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			panic(err)
		}

		defer file.Close()

		log.SetOutput(file)
	}

	raven.CaptureError(fmt.Errorf(message, v...), nil)
	log.Printf(message, v...)

	if l.Debug {
		fmt.Printf("[ERROR] "+message+"\n", v...)
	}
}

// Info logs everything that happens in application
func (l Logs) Info(message string, v ...interface{}) {
	if len(l.InfoLog) > 0 {
		file, err := os.OpenFile(l.InfoLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		log.SetOutput(file)
	}

	log.Printf(message, v...)

	if l.Debug {
		fmt.Printf("[INFO] "+message+"\n", v...)
	}
}

// Logger holds methods to log messages to files
var Logger = new(Logs)
