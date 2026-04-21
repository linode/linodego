package linodego

import (
	"log"
	"os"
)

type Logger interface {
	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
}

type logger struct {
	l *log.Logger
}

func createLogger() *logger {
	l := &logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

var _ Logger = (*logger)(nil)

func (l *logger) Errorf(format string, v ...any) {
	l.output("ERROR "+format, v...)
}

func (l *logger) Warnf(format string, v ...any) {
	l.output("WARN "+format, v...)
}

func (l *logger) Debugf(format string, v ...any) {
	l.output("DEBUG "+format, v...)
}

func (l *logger) output(format string, v ...any) { //nolint:goprintffuncname
	if len(v) == 0 {
		l.l.Print(format)
		return
	}

	l.l.Printf(format, v...)
}
