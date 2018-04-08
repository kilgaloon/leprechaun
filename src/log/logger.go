package log

// Logger interface holds methods to log messages
type Logger interface {
	error(message string, v ... interface{})
	access(message string, v ... interface{})
}

// FlushError to log
func FlushError(l Logger, message string, v ... interface{}) {
	l.error(message, v ...)
}

// FlushAccess to log
func FlushAccess(l Logger, message string, v ... interface{}) {
	l.access(message, v ...)
}
