package log

import (
	"fmt"

	"github.com/degatedev/degatesdk/conf"
)

func Print(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Info(format string, args ...interface{}) {
	if conf.Debug {
		fmt.Printf(format+"\n", args...)
	}
}

func Error(format string, args ...interface{}) {
	if conf.Debug {
		fmt.Printf(format+"\n", args...)
	}
}
