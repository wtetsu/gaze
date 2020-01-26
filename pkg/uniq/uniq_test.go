/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package uniq

import (
	"testing"
)

func Test(t *testing.T) {
	uniq := New()
	if len(uniq.List()) != 0 {
		t.Fatal()
	}

	uniq.Add("aaa")
	uniq.Add("bbb")

	if len(uniq.List()) != 2 {
		t.Fatal()
	}

	uniq.Add("bbb")
	uniq.Add("bbb")
	uniq.Add("ccc")
	uniq.Add("ccc")

	if uniq.Len() != 3 {
		t.Fatal()
	}

	uniq.AddAll([]string{"bbb", "ccc", "ddd"})
	if uniq.Len() != 4 {
		t.Fatal()
	}
}
