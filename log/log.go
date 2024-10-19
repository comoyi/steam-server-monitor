package log

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"time"
)

var logger Logger

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
	LevelFatal Level = 600
	LevelOff   Level = 900
)

const (
	NameTrace string = "TRACE"
	NameDebug string = "DEBUG"
	NameInfo  string = "INFO"
	NameWarn  string = "WARN"
	NameError string = "ERROR"
	NameFatal string = "FATAL"
	NameOff   string = "OFF"
)

var levelNameMap = map[Level]string{
	LevelTrace: NameTrace,
	LevelDebug: NameDebug,
	LevelInfo:  NameInfo,
	LevelWarn:  NameWarn,
	LevelError: NameError,
	LevelFatal: NameFatal,
	LevelOff:   NameOff,
}

var nameLevelMap = map[string]Level{
	NameTrace: LevelTrace,
	NameDebug: LevelDebug,
	NameInfo:  LevelInfo,
	NameWarn:  LevelWarn,
	NameError: LevelError,
	NameFatal: LevelFatal,
	NameOff:   LevelOff,
}

func GetLevel() Level {
	return logger.GetLevel()
}

func SetLevel(level Level) Level {
	return logger.SetLevel(level)
}

func SetLevelByName(levelName string) (Level, error) {
	level, err := GetLevelByName(levelName)
	if err != nil {
		return 0, err
	}
	return SetLevel(level), nil
}

func GetLevelByName(levelName string) (Level, error) {
	name := strings.ToUpper(levelName)
	if level, ok := nameLevelMap[name]; ok {
		return level, nil
	}
	return 0, errors.New("invalid level name")
}

func Tracef(format string, args ...any) {
	logger.Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func Init() error {
	logger = New()
	return nil
}

func New() Logger {
	logger := &MyLogger{
		level: LevelInfo,
	}
	return logger
}

type Logger interface {
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	GetLevel() Level
	SetLevel(level Level) Level
	logf(level Level, format string, args ...any)
}

type MyLogger struct {
	level Level
}

func (l *MyLogger) Tracef(format string, args ...any) {
	l.logf(LevelTrace, format, args...)
}

func (l *MyLogger) Debugf(format string, args ...any) {
	l.logf(LevelDebug, format, args...)
}

func (l *MyLogger) Infof(format string, args ...any) {
	l.logf(LevelInfo, format, args...)
}

func (l *MyLogger) Warnf(format string, args ...any) {
	l.logf(LevelWarn, format, args...)
}

func (l *MyLogger) Errorf(format string, args ...any) {
	l.logf(LevelError, format, args...)
}

func (l *MyLogger) Fatalf(format string, args ...any) {
	l.logf(LevelFatal, format, args...)
}

func (l *MyLogger) GetLevel() Level {
	return l.level
}

func (l *MyLogger) SetLevel(level Level) Level {
	oldLevel := l.level
	l.level = level
	return oldLevel
}

func (l *MyLogger) logf(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	_, filepath, line, _ := runtime.Caller(2)

	writeToConsole(level, filepath, line, format, args...)
	writeToFile(level, filepath, line, format, args...)
}

func writeToConsole(level Level, filepath string, line int, format string, args ...any) {
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
	fileLine := fmt.Sprintf("%s:%d", filepath, line)
	levelText := fmt.Sprintf("%s%s%s%s%s%s", color, "[", level.String(), "]", colorReset, postfix)
	s := fmt.Sprintf(fmt.Sprintf("%s\n%s%s", fileLine, time.Now().Format(time.DateTime), levelText)+format, args...)
	io.WriteString(os.Stdout, s+"\n")
}

func writeToFile(level Level, filepath string, line int, format string, args ...any) {
	fileLine := fmt.Sprintf("%s:%d", filepath, line)
	s := fmt.Sprintf(fmt.Sprintf("%s\n%s%-7s", fileLine, time.Now().Format(time.DateTime), "["+level.String()+"]")+format, args...)
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
