/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package tutil

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	zero := GetFileModifiedTime("___invalid__")
	if zero != 0 {
		t.Fatal()
	}

	fileTime := GetFileModifiedTime("tutil.go")
	if fileTime == 0 {
		t.Fatal()
	}

	time.Sleep(1 * time.Millisecond)
	now1 := time.Now().UnixNano()

	if now1 < fileTime {
		t.Fatal()
	}

	ch := After(5)
	now2 := time.Now().UnixNano()
	if now2 < now1 {
		t.Fatal()
	}
	<-ch
	now3 := time.Now().UnixNano()
	if now3 < now2 {
		t.Fatal()
	}
}
