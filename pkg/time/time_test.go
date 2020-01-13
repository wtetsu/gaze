/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package time

import (
	"testing"
)

func Test(t *testing.T) {
	zero := GetFileModifiedTime("___invalid__")
	if zero != 0 {
		t.Fatal()
	}

	fileTime := GetFileModifiedTime("time.go")

	now1 := Now()

	if now1 < fileTime {
		t.Fatal()
	}

	ch := After(5)
	now2 := Now()
	if now2 < now1 {
		t.Fatal()
	}
	<-ch
	now3 := Now()
	if now3 < now2 {
		t.Fatal()
	}
}
