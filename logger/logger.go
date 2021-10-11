package logger

import (
	"io"
	"log"
)

// Logger represents logger interface.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// CustomLog represents logger struct.
type CustomLog struct {
	defaultLogger *log.Logger
}

// New returns logger object.
func New(out io.Writer) *CustomLog {
	l := new(CustomLog)
	l.defaultLogger = log.New(out, "", log.Ldate|log.Ltime)

	return l
}

// Info used for info messages.
func (l *CustomLog) Info(v ...interface{}) {
	l.defaultLogger.Println(v...)
}

// Infof used for formatted info messages.
func (l *CustomLog) Infof(format string, v ...interface{}) {
	l.defaultLogger.Printf("INFO: "+format, v...)
}

// Error used for error messages.
func (l *CustomLog) Error(v ...interface{}) {
	l.defaultLogger.Println(v...)
}

// Errorf used for formatted error messages.
func (l *CustomLog) Errorf(format string, v ...interface{}) {
	l.defaultLogger.Printf("ERROR: "+format, v...)
}

// Debug used for debug messages.
func (l *CustomLog) Debug(v ...interface{}) {
	l.defaultLogger.Print(v...)
}

// Debugf used for formatted debug message.
func (l *CustomLog) Debugf(format string, v ...interface{}) {
	l.defaultLogger.Printf("DEBUG: "+format, v...)
}

// Warning used for warning message.
func (l *CustomLog) Warning(v ...interface{}) {
	l.defaultLogger.Print(v...)
}

// Warningf used for formatted warning message.
func (l *CustomLog) Warningf(format string, v ...interface{}) {
	l.defaultLogger.Printf("WARNING: "+format, v...)
}
