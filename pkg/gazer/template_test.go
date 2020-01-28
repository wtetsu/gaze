/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"testing"
)

func TestTemplate1(t *testing.T) {
	r, err := render("{{file}} {{ext}} {{base}} {{dir}} {{base0}} {{base1}} {{base2}}", "/full/path/test.txt.bak")
	if err != nil {
		t.Fatal(err)
	}

	if r != "/full/path/test.txt.bak .bak test.txt.bak /full/path test test.txt test.txt.bak" {
		t.Fatal(r)
	}
}

func TestTemplateError(t *testing.T) {
	r, err := render("{{file}", "/full/path/test.txt.bak")
	if err == nil || r != "" {
		t.Fatal(err)
	}
}
