package log

import (
	"log"
	"os"
)

// Logs struct holds path to different loggs
type Logs struct {
	errorLog string
}

func (l Logs) Error(message string, v ...interface{}) {
	file, err := os.OpenFile(l.errorLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Fatalf(message, v...)
}

// CreateLogs and return struct
func CreateLogs(errorLog string) Logs {
	return Logs{errorLog}
}
