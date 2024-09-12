package paint

import "fmt"

func InfoF(format string, a ...interface{}) {
	info.Printf(format, a...)
	fmt.Println()
}

func ErrorF(format string, a ...interface{}) {
	err.Printf(format, a...)
	fmt.Println()
}

func WarnF(format string, a ...interface{}) {
	warn.Printf(format, a...)
	fmt.Println()
}

func SuccessF(format string, a ...interface{}) {
	success.Printf(format, a...)
	fmt.Println()
}
