package debuglog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// findFunc formats the function name and line number for log output.
func findFunc(f *runtime.Frame) (string, string) {
	s := strings.Split(f.Function, ".")
	funcname := s[len(s)-1]
	return fmt.Sprintf("%s:%d", funcname, f.Line), ""
}

// Field sequence for custom sorting
var fieldSeq = map[string]int{
	"time":  0,
	"level": 1,
	"func":  2,
}

// sortCustom customizes the sorting of log fields.
func sortCustom(fields []string) {
	sort.Slice(fields, func(i, j int) bool {
		iIdx, oki := fieldSeq[fields[i]]
		jIdx, okj := fieldSeq[fields[j]]
		if oki && okj {
			return iIdx < jIdx
		}
		if oki {
			return true
		}
		if okj {
			return false
		}
		return fields[i] < fields[j] // Alphabetical order for new fields
	})
}

// DebugLogInit initializes a logger with flexible options.
func DebugLogInit(logname string, makedir bool, usePID bool, useMultiWriter bool) *logrus.Logger {
	debuglogrus := logrus.New()

	// Set log formatter
	debugFormatter := &logrus.TextFormatter{
		DisableColors:    true,
		CallerPrettyfier: findFunc,
		ForceQuote:       true,
		DisableSorting:   false,
		SortingFunc:      sortCustom,
	}
	debuglogrus.SetReportCaller(true)
	debuglogrus.SetFormatter(debugFormatter)

	var logpath string

	baselogpath := os.Getenv("LOG_BASE_PATH")
	if baselogpath == "" {
		var err error
		baselogpath, err = os.Getwd()
		if err != nil {
			fmt.Println("Cannot get Current Directory")
			baselogpath = "./"
		}
	}

	if makedir {
		makepath := fmt.Sprintf("%s/log", baselogpath)
		if err := os.MkdirAll(makepath, os.ModePerm); err != nil {
			fmt.Printf("Cannot make log Directory: %s\n", err.Error())
			fmt.Println("Logging will proceed in the current directory.")
			logpath = baselogpath
		} else {
			logpath = makepath
		}
	} else {
		logpath = baselogpath
	}

	var debuglogpath string
	if usePID {
		debuglogpath = fmt.Sprintf("%s/%s.%d.log", logpath, logname, os.Getpid())
	} else {
		debuglogpath = fmt.Sprintf("%s/%s.%s.log", logpath, logname, time.Now().Format("20060102_150405"))
	}

	debugLogOutput := &lumberjack.Logger{
		Filename:   debuglogpath,
		MaxSize:    500,  // Max file size in MB
		MaxBackups: 3,    // Max backup files
		MaxAge:     3,    // Max age in days
		Compress:   true, // Enable compression
	}

	if useMultiWriter {
		multiWriter := io.MultiWriter(debugLogOutput, os.Stdout)
		debuglogrus.SetOutput(multiWriter)
	} else {
		debuglogrus.SetOutput(debugLogOutput)
	}

	fmt.Printf("Logging to %s\n", debuglogpath)
	return debuglogrus
}
