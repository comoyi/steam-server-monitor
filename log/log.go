package log

import (
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"os"
)

func Debugf(format string, args ...interface{}) {
	//fmt.Printf("[DEBUG] "+format, args...)
	s := fmt.Sprintf("[DEBUG]"+format, args...)
	w(s)
}

func Infof(format string, args ...interface{}) {
	//fmt.Printf("[INFO] "+format, args...)
	s := fmt.Sprintf("[INFO] "+format, args...)
	w(s)
}

func Warnf(format string, args ...interface{}) {
	//fmt.Printf("[WARN] "+format, args...)
	s := fmt.Sprintf("[WARN] "+format, args...)
	w(s)
}

func Errorf(format string, args ...interface{}) {
	//fmt.Printf("[ERROR]"+format, args...)
	s := fmt.Sprintf("[ERROR]"+format, args...)
	w(s)
}

func w(s string) {
	if !viper.GetBool("debug") {
		return
	}
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.WriteString(s)
	if err != nil {
		return
	}
}
