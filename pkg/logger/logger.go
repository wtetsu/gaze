package logger

import "fmt"

// Log level.
const (
	SILENT  = 0
	QUIET   = 1
	NORMAL  = 2
	VERBOSE = 3
	DEBUG   = 4
)

var logLevel int = NORMAL

// Level sets a new log level.
func Level(newLogLevel int) {
	logLevel = newLogLevel
}

// Error writes a fatal log
func Error(a ...interface{}) {
	if logLevel < QUIET {
		return
	}

	fmt.Println(a...)
}

// Notice writes a warning log
func Notice(format string, a ...interface{}) {
	if logLevel <= NORMAL {
		return
	}
	fmt.Printf(format, a...)
	fmt.Println()
}

// Info writes a info log
func Info(format string, a ...interface{}) {
	if logLevel < VERBOSE {
		return
	}
	fmt.Printf(format, a...)
	fmt.Println()
}

// Debug writes a debug log
func Debug(format string, a ...interface{}) {
	if logLevel < DEBUG {
		return
	}
	fmt.Printf(format, a...)
	fmt.Println()
}

// DebugObject writes a debug log
func DebugObject(a ...interface{}) {
	if logLevel < DEBUG {
		return
	}
	fmt.Println(a...)
}
