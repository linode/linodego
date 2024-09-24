package linodego

import (
	"log"
	"os"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type logger struct {
	l *log.Logger
}

// nolint: unused
func createLogger() *logger {
	l := &logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

var _ Logger = (*logger)(nil)

func (l *logger) Errorf(format string, v ...interface{}) {
	l.output("ERROR "+format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.output("WARN "+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.output("DEBUG "+format, v...)
}

func (l *logger) output(format string, v ...interface{}) { //nolint:goprintffuncname
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
