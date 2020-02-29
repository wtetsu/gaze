/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package time

import (
	"os"
	"time"
)

// GetFileModifiedTime returns the last modified time of a file.
// If there is an error, it will return 0
func GetFileModifiedTime(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.ModTime().UnixNano()
}

// Now returns current time as a int64 value.
func Now() int64 {
	return time.Now().UnixNano()
}

// Sleep sleeps milliseconds time
func Sleep(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}

// After waits for the duration.
func After(d int64) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(d) * time.Millisecond)
		close(ch)
	}()

	return ch
}
