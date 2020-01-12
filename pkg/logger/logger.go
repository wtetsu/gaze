package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Log level.
const (
	SILENT  = 0
	QUIET   = 1
	NORMAL  = 2
	VERBOSE = 3
	DEBUG   = 4
)

var logLevel = NORMAL
var count = 0

var printInfo func(format string, a ...interface{})
var printNotice func(format string, a ...interface{})
var printError func(format string, a ...interface{})

// Level sets a new log level.
func Level(newLogLevel int) {
	logLevel = newLogLevel
}

// Colorful enables colorful output
func Colorful() {
	printInfo = color.New(color.FgHiCyan).PrintfFunc()
	printNotice = color.New(color.FgCyan).PrintfFunc()

	f := color.New(color.FgRed).FprintfFunc()
	printError = func(format string, a ...interface{}) {
		f(color.Error, format, a...)
	}
}

// Plain disables colorful output
func Plain() {
	printInfo = func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
	printNotice = func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
	printError = func(format string, a ...interface{}) {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// Error writes an error log
func Error(format string, a ...interface{}) {
	if logLevel < QUIET {
		return
	}
	space()
	printError(format, a...)
	fmt.Println(a...)
	fmt.Println()
	count++
}

// ErrorObject writes an error log
func ErrorObject(a ...interface{}) {
	Error("%v", a...)
}

// Notice writes a notice log
func Notice(format string, a ...interface{}) {
	notice(false, format, a...)
}

// NoticeWithBlank writes a notice log
func NoticeWithBlank(format string, a ...interface{}) {
	notice(true, format, a...)
}

// NoticeObject writes a notice log
func NoticeObject(a ...interface{}) {
	notice(false, "%v", a...)
}

func notice(enableSpace bool, format string, a ...interface{}) {
	if logLevel < NORMAL {
		return
	}
	if enableSpace {
		space()
	}
	printNotice(format, a...)
	fmt.Println()
	count++
}

// Info writes a info log
func Info(format string, a ...interface{}) {
	if logLevel < VERBOSE {
		return
	}
	printInfo(format, a...)
	fmt.Println()
	count++
}

// Debug writes a debug log
func Debug(format string, a ...interface{}) {
	if logLevel < DEBUG {
		return
	}
	fmt.Printf(format, a...)
	fmt.Println()
	count++
}

// DebugObject writes a debug log
func DebugObject(a ...interface{}) {
	Debug("%v", a...)
}

func space() {
	count++
	if count <= 1 {
		return
	}
	fmt.Println()
}
