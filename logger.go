package migorm

import "fmt"

const (
	InfoLevel  = "INFO"
	ErrorLevel = "ERROR"
)

type Logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

func NewLogger() Logger {
	return &logger{}
}

type logger struct {
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.log(InfoLevel, template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.log(ErrorLevel, template, args...)
}

func (logger) log(level string, template string, args ...interface{}) {
	msg := template
	if msg == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(template, args...)
	}

	fmt.Printf(level+": %v\n", msg)
}
