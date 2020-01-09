package logger

import "fmt"

// Debug writes a debug log
func Debug(a ...interface{}) {
	fmt.Println(a...)
}

// Debugf writes a debug log
func Debugf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
}

// Info writes a info log
func Info(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
}

// Fatal writes a fatal log
func Fatal(a ...interface{}) {
	fmt.Println(a...)
}
