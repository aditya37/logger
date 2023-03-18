package logger

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

var logrusSingleton sync.Once

// logrus instance
var l *logrus.Logger

// logrus entry
var e *logrus.Entry

func init() {
	logrusSingleton.Do(func() {
		l = logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: false,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				s := strings.Split(f.Function, ".")
				funcname := s[len(s)-1]
				return funcname, f.File + ":" + strconv.Itoa(f.Line)
			},
		}
		// output
		l.SetOutput(os.Stdout)
		// apm hook
		l.AddHook(&apmlogrus.Hook{})
		// logrus entry
		e = logrus.NewEntry(l)
	})
}

// set log level
func SetLevel(level logrus.Level) {
	l.SetLevel(level)
	e = logrus.NewEntry(l)
}

// set logrus or log type
func Info(args ...interface{})  { e.Info(args...) }
func Debug(args ...interface{}) { e.Debug(args...) }
func Warn(args ...interface{})  { e.Warn(args...) }
func Error(args ...interface{}) { e.Error(args...) }
