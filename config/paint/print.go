package paint

import "github.com/fatih/color"

func log(logger *color.Color, a ...interface{}) {
	logger.Println(a...)
}

func Info(a ...interface{}) {
	log(info, a...)
}

func Error(a ...interface{}) {
	log(err, a...)
}

func Warn(a ...interface{}) {
	log(warn, a...)
}

func Success(a ...interface{}) {
	log(success, a...)
}

func Notice(a ...interface{}) {
	log(notice, a...)
}
