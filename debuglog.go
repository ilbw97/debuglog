package debuglog

import (
	"errors"
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
	LogName         string `json:"name"`
	MakeDir         bool   `json:"make_dir" example:"true"`
	UsePID          bool   `json:"use_pid" example:"true"`
	UseMultiWriter  bool   `json:"use_multi_writer" example:"false"`
	LogRotateConfig `json:"rotate_config"`
}

type LogRotateConfig struct {
	MaxSize    int  `json:"max_size" example:"500"`  // 최대 로그 파일 크기 (MB)
	MaxBackups int  `json:"max_backups" example:"3"` // 최대 백업 파일 개수
	MaxAge     int  `json:"max_age" example:"3"`     // 로그 파일 최대 보관 기간 (일)
	Compress   bool `json:"compress" example:"true"` // 오래된 로그 파일 압축 여부
}

// 기본값 설정 함수
func setDefaultLogConfig(config *LogConfig) {
	if config.MaxSize == 0 {
		config.MaxSize = 500
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 3
	}
	if config.MaxAge == 0 {
		config.MaxAge = 3
	}
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

// LogInit initializes a logger with flexible options.
func LogInit(config *LogConfig) (*logrus.Logger, error) {
	if config == nil {
		errMsg := errors.New("config is nil")
		logrus.Error(errMsg)
		return nil, errMsg
	}

	setDefaultLogConfig(config)

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

	logpath, err := determineLogPath(config)
	if err != nil {
		logrus.Warnf("Falling back to current directory for logging: %v", err)
		logpath = "./"
	}

	var debuglogpath string
	if config.UsePID {
		debuglogpath = fmt.Sprintf("%s/%s.%d.log", logpath, config.LogName, os.Getpid())
	} else {
		debuglogpath = fmt.Sprintf("%s/%s.%s.log", logpath, config.LogName, time.Now().Format("20060102_150405"))
	}

	debugLogOutput := &lumberjack.Logger{
		Filename:   debuglogpath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	if config.UseMultiWriter {
		multiWriter := io.MultiWriter(debugLogOutput, os.Stdout)
		debuglogrus.SetOutput(multiWriter)
	} else {
		debuglogrus.SetOutput(debugLogOutput)
	}

	logrus.Infof("Logging initialized at %s", debuglogpath)
	return debuglogrus, nil
}

// determineLogPath determines the log directory path based on configuration.
func determineLogPath(config *LogConfig) (string, error) {
	baselogpath := os.Getenv("LOG_BASE_PATH")
	if baselogpath == "" {
		var err error
		baselogpath, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("cannot get current directory: %w", err)
		}

		logrus.Infof("LOG_BASE_PATH not set. Using working directory: %s", baselogpath)
	}

	if config.MakeDir {
		makepath := fmt.Sprintf("%s/log", baselogpath)
		if err := os.MkdirAll(makepath, os.FileMode(os.O_RDWR)); err != nil {
			return "", fmt.Errorf("cannot create log directory: %w", err)
		}
		return makepath, nil
	}

	return baselogpath, nil
}
