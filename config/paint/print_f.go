package paint

import (
	"fmt"

	"github.com/fatih/color"
)

func logF(logger *color.Color, format string, a ...interface{}) {
	logger.Printf(format, a...)
	fmt.Println()
}

func InfoF(format string, a ...interface{}) {
	logF(info, format, a...)
}

func ErrorF(format string, a ...interface{}) {
	logF(err, format, a...)
}

func WarnF(format string, a ...interface{}) {
	logF(warn, format, a...)
}

func SuccessF(format string, a ...interface{}) {
	logF(success, format, a...)
}

func NoticeF(format string, a ...interface{}) {
	logF(notice, format, a...)
}
