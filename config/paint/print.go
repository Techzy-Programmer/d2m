package paint

func Info(a ...interface{}) {
	info.Println(a...)
}

func Error(a ...interface{}) {
	err.Println(a...)
}

func Warn(a ...interface{}) {
	warn.Println(a...)
}

func Success(a ...interface{}) {
	success.Println(a...)
}
