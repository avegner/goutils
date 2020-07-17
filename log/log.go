package log

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type Level int

const (
	CriticalLevel Level = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
)

var levelPrefixes = []string{"CRI", "ERR", "WRN", "INF", "DBG"}

type Logger interface {
	SetLevel(l Level)
	Level() Level
	Criticalf(format string, args ...interface{}) (n int, err error)
	Errorf(format string, args ...interface{}) (n int, err error)
	Warningf(format string, args ...interface{}) (n int, err error)
	Infof(format string, args ...interface{}) (n int, err error)
	Debugf(format string, args ...interface{}) (n int, err error)
}

type logger struct {
	w   io.Writer
	l   Level
	lmu sync.RWMutex
	m   string
}

func New(w io.Writer, l Level, m string) Logger {
	return &logger{
		w: w,
		l: l,
		m: normalize(m),
	}
}

func (log *logger) SetLevel(l Level) {
	log.lmu.Lock()
	defer log.lmu.Unlock()
	log.l = l
}

func (log *logger) Level() Level {
	return log.level()
}

func (log *logger) Criticalf(format string, args ...interface{}) (n int, err error) {
	return log.printf(CriticalLevel, format, args...)
}

func (log *logger) Errorf(format string, args ...interface{}) (n int, err error) {
	return log.printf(ErrorLevel, format, args...)
}

func (log *logger) Warningf(format string, args ...interface{}) (n int, err error) {
	return log.printf(WarningLevel, format, args...)
}
func (log *logger) Infof(format string, args ...interface{}) (n int, err error) {
	return log.printf(InfoLevel, format, args...)
}

func (log *logger) Debugf(format string, args ...interface{}) (n int, err error) {
	return log.printf(DebugLevel, format, args...)
}

func (log *logger) printf(l Level, format string, args ...interface{}) (n int, err error) {
	if l > log.level() {
		return
	}
	return log.w.Write(
		[]byte(fmt.Sprintf("[%s][%s] ", levelPrefixes[l], log.m) + fmt.Sprintf(format, args...)))
}

func (log *logger) level() Level {
	log.lmu.RLock()
	defer log.lmu.RUnlock()
	return log.l
}

func normalize(prefix string) string {
	return strings.ToUpper(prefix)
}
