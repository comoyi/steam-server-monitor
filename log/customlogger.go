package log

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"time"
)

func New() Logger {
	logger := &CustomLogger{
		level: LevelInfo,
	}
	return logger
}

type CustomLogger struct {
	level Level
}

func (l *CustomLogger) Tracef(format string, args ...any) {
	l.logf(LevelTrace, format, args...)
}

func (l *CustomLogger) Debugf(format string, args ...any) {
	l.logf(LevelDebug, format, args...)
}

func (l *CustomLogger) Infof(format string, args ...any) {
	l.logf(LevelInfo, format, args...)
}

func (l *CustomLogger) Warnf(format string, args ...any) {
	l.logf(LevelWarn, format, args...)
}

func (l *CustomLogger) Errorf(format string, args ...any) {
	l.logf(LevelError, format, args...)
}

func (l *CustomLogger) Fatalf(format string, args ...any) {
	l.logf(LevelFatal, format, args...)
}

func (l *CustomLogger) Enabled(level Level) bool {
	if level >= l.level {
		return true
	}
	return false
}

func (l *CustomLogger) GetLevel() Level {
	return l.level
}

func (l *CustomLogger) SetLevel(level Level) Level {
	oldLevel := l.level
	l.level = level
	return oldLevel
}

func (l *CustomLogger) logf(level Level, format string, args ...any) {
	if !l.Enabled(level) {
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
