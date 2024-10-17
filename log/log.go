package log

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
)

var logger *Logger

var currentLogLevel Level = LevelOff

type Level int

func (l Level) String() string {
	if name, ok := levelNameMap[l]; ok {
		return name
	}
	return ""
}

const (
	LevelTrace Level = 100
	LevelDebug Level = 200
	LevelInfo  Level = 300
	LevelWarn  Level = 400
	LevelError Level = 500
	LevelOff   Level = 900
)

const (
	NameTrace string = "TRACE"
	NameDebug string = "DEBUG"
	NameInfo  string = "INFO"
	NameWarn  string = "WARN"
	NameError string = "ERROR"
	NameOff   string = "OFF"
)

var levelNameMap = map[Level]string{
	LevelTrace: NameTrace,
	LevelDebug: NameDebug,
	LevelInfo:  NameInfo,
	LevelWarn:  NameWarn,
	LevelError: NameError,
	LevelOff:   NameOff,
}

var nameLevelMap = map[string]Level{
	NameTrace: LevelTrace,
	NameDebug: LevelDebug,
	NameInfo:  LevelInfo,
	NameWarn:  LevelWarn,
	NameError: LevelError,
	NameOff:   LevelOff,
}

func LogLevel() Level {
	return currentLogLevel
}

func SetLogLevel(level Level) Level {
	oldLogLevel := currentLogLevel
	currentLogLevel = level
	return oldLogLevel
}

func SetLogLevelByName(levelName string) (Level, error) {
	level, err := GetLogLevelByName(levelName)
	if err != nil {
		return 0, err
	}
	return SetLogLevel(level), nil
}

func GetLogLevelByName(levelName string) (Level, error) {
	name := strings.ToUpper(levelName)
	if level, ok := nameLevelMap[name]; ok {
		return level, nil
	}
	return 0, errors.New("invalid level name")
}

func Tracef(format string, args ...any) {
	logger.logf(LevelTrace, format, args...)
}

func Debugf(format string, args ...any) {
	logger.logf(LevelDebug, format, args...)
}

func Infof(format string, args ...any) {
	logger.logf(LevelInfo, format, args...)
}

func Warnf(format string, args ...any) {
	logger.logf(LevelWarn, format, args...)
}

func Errorf(format string, args ...any) {
	logger.logf(LevelError, format, args...)
}

func Init() error {
	logger = New()
	return nil
}

func New() *Logger {
	logger := &Logger{}
	return logger
}

type Logger struct{}

func (l *Logger) logf(level Level, format string, args ...any) {
	if level < currentLogLevel {
		return
	}

	writeToConsole(level, format, args...)
	writeToFile(level, format, args...)
}

func writeToConsole(level Level, format string, args ...any) {
	color := ""
	postfix := ""
	switch level {
	case LevelTrace:
		color = ""
	case LevelDebug:
		color = colorGreen
	case LevelInfo:
		color = colorBlue
		postfix = " "
	case LevelWarn:
		color = colorYellow
		postfix = " "
	case LevelError:
		color = colorRed
	}
	levelText := fmt.Sprintf("%s%s%s%s%s%s", color, "[", level.String(), "]", colorReset, postfix)
	s := fmt.Sprintf(fmt.Sprintf("%s%s", time.Now().Format(time.DateTime), levelText)+format, args...)
	io.WriteString(os.Stdout, s+"\n")
}

func writeToFile(level Level, format string, args ...any) {
	s := fmt.Sprintf(fmt.Sprintf("%s%-7s", time.Now().Format(time.DateTime), "["+level.String()+"]")+format, args...)
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.WriteString(s + "\n")
	if err != nil {
		return
	}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[1;34m"
)
