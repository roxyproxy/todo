package logger

import (
	"io"
	"log"
)

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

type CustomLog struct {
	defaultLogger *log.Logger
}

func New(out io.Writer) *CustomLog {
	l := new(CustomLog)
	l.defaultLogger = log.New(out, "", log.Ldate|log.Ltime)

	return l
}

func (l *CustomLog) Info(v ...interface{}) {
	l.defaultLogger.Println(v)
}

func (l *CustomLog) Infof(format string, v ...interface{}) {
	l.defaultLogger.Printf("INFO: "+format, v...)
}

func (l *CustomLog) Error(v ...interface{}) {
	l.defaultLogger.Println(v)
}

func (l *CustomLog) Errorf(format string, v ...interface{}) {
	l.defaultLogger.Printf("ERROR: "+format, v...)
}

func (l *CustomLog) Debug(v ...interface{}) {
	l.defaultLogger.Print(v)
}

func (l *CustomLog) Debugf(format string, v ...interface{}) {
	l.defaultLogger.Printf("DEBUG: "+format, v...)
}

func (l *CustomLog) Warning(v ...interface{}) {
	l.defaultLogger.Print(v)
}

func (l *CustomLog) Warningf(format string, v ...interface{}) {
	l.defaultLogger.Printf("WARNING: "+format, v...)
}
