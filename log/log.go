package log

import "fmt"

func Debugf(format string, args ...interface{}) {
	fmt.Printf("[DEBUG]"+format, args...)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format, args...)
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format, args...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR]"+format, args...)
}
