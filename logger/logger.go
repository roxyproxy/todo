package logger

import "log"

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type defaultLogger struct {
	l *log.Logger
}

/*
func New() *defaultLogger {
	return &defaultLogger{l: log.New()}
}
*/
