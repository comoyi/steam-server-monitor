package log

import (
	"errors"
	"strings"
)

var defaultLogger Logger

func Init() error {
	logger := New()
	SetDefault(logger)
	return nil
}

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

func Default() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
}

func GetLevel() Level {
	return Default().GetLevel()
}

func SetLevel(level Level) Level {
	return Default().SetLevel(level)
}

func GetLevelByName(levelName string) (Level, error) {
	name := strings.ToUpper(levelName)
	if level, ok := nameLevelMap[name]; ok {
		return level, nil
	}
	return 0, errors.New("invalid level name")
}

func SetLevelByName(levelName string) (Level, error) {
	level, err := GetLevelByName(levelName)
	if err != nil {
		return 0, err
	}
	return SetLevel(level), nil
}

func Tracef(format string, args ...any) {
	Default().Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}

func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	Default().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	Default().Fatalf(format, args...)
}
