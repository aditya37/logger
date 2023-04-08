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
var (

	// qualified package name, cached at first use
	logrusPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// logrus entry
var e *logrus.Entry

func init() {
	logrusSingleton.Do(func() {
		l = logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: false,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				f = getCaller()
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

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != "github.com/sirupsen/logrus" {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
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
