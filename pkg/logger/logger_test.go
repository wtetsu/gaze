/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package logger

import (
	"testing"
)

func Test(t *testing.T) {
	Colorful()
	for level := 0; level <= 4; level++ {
		Level(level)
		writeAll()
	}

	Plain()
	for level := 0; level <= 4; level++ {
		Level(level)
		Colorful()
	}
}

func writeAll() {
	Error("log")
	ErrorObject("log")

	Notice("log")
	NoticeObject("log")
	NoticeWithBlank("log")

	Info("log")
	//InfoObject("log")

	Debug("log")
	DebugObject("log")
}
