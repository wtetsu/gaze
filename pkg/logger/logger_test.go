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
	Error("log(Error)")
	ErrorObject("log(ErrorObject)")

	Notice("log(Notice)")
	NoticeObject("log(NoticeObject)")
	NoticeWithBlank("log(NoticeWithBlank)")

	Info("log(Info)")
	//InfoObject("log")

	Debug("log(Debug)")
	DebugObject("log(DebugObject)")
}
