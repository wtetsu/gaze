package logger

import "fmt"

func Debug(a ...interface{}) {
	fmt.Println(a...)
}

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Fatal(a ...interface{}) {
	fmt.Println(a...)
}
