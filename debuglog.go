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

type LogConfig struct {
	LogName     string `json:"name"`
	MakeDir     bool	`json:"make_dir" example:"true"`
	UsePID      bool	`json:"use_pid" exmaple:"true"`
	UseMultiWriter bool	`json:"use_multi_writer" example:"false"`
	LogRotateConfig
}

type LogRotateConfig struct {
	MaxSize int `json:"max_size" exmaple:"500"`
	MaxBackups int `json:"max_backups" example:"3"`
	MaxAge int `json:"max_age" example:"3"`
	Compress bool `json:"compress" exmaple:"true"`
}

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
func DebugLogInit(options *LogOptions) *logrus.Logger {
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

	if options.MakeDir {
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
	if options.UsePID {
		debuglogpath = fmt.Sprintf("%s/%s.%d.log", logpath, options.LogName, os.Getpid())
	} else {
		debuglogpath = fmt.Sprintf("%s/%s.%s.log", logpath, options.LogName, time.Now().Format("20060102_150405"))
	}

	debugLogOutput := &lumberjack.Logger{
		Filename:   debuglogpath,
		MaxSize:    options.MaxSize,  // Max file size in MB
		MaxBackups: options.MaxBackups,    // Max backup files
		MaxAge:     options.MaxAge,    // Max age in days
		Compress:   options.Compress, // Enable compression
	}

	if options.UseMultiWriter {
		multiWriter := io.MultiWriter(debugLogOutput, os.Stdout)
		debuglogrus.SetOutput(multiWriter)
	} else {
		debuglogrus.SetOutput(debugLogOutput)
	}

	fmt.Printf("Logging to %s\n", debuglogpath)

	return debuglogrus
}
