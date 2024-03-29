package debuglog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func findFunc(f *runtime.Frame) (string, string) {
	s := strings.Split(f.Function, ".")
	funcname := s[len(s)-1]

	return funcname, ""
}

var fieldSeq = map[string]int{
	"time":  0,
	"level": 1,
	"func":  2,
}

func sortCustom(fields []string) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i] == "msg" {
			return false
		}
		if iIdx, oki := fieldSeq[fields[i]]; oki {
			if jIdx, okj := fieldSeq[fields[j]]; okj {
				return iIdx < jIdx
			}
			return true
		}
		return false
	})
}

func DebugLogInit(logname string, makedir bool) *logrus.Logger {

	debuglogrus := logrus.New()
	debugFormatter := &logrus.TextFormatter{
		DisableColors:    true,
		CallerPrettyfier: findFunc,
		ForceQuote:       true,
		DisableSorting:   false,
		SortingFunc:      sortCustom,
	}

	debuglogrus.SetReportCaller(true)
	debuglogrus.SetFormatter(debugFormatter)

	// path
	var logpath string

	baselogpath, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get CurrentDirectory")
	}
	fmt.Println("baselogpath :" + baselogpath)
	if len(baselogpath) == 0 {
		baselogpath = "./"
	}

	var makepath string
	if !makedir {
		logpath = baselogpath
	} else {
		makepath = baselogpath + "/log"
		fmt.Println("makepath :" + makepath)
		if _, err := os.Stat(makepath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(makepath, os.ModePerm)
			if err != nil {
				fmt.Println("Cannot make log Directory" + err.Error())
				fmt.Println("Trying to logging at Current Directory")
			}
		}
		logpath = makepath
	}

	// WAF LOG OUTPUT SETTING
	debuglogpath := fmt.Sprintf("%s/%s.%d.log", logpath, logname, os.Getpid())
	debugLogOutput := lumberjack.Logger{
		Filename:   debuglogpath,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   false,
	}

	debuglogrus.SetOutput(&debugLogOutput)
	return debuglogrus
}
