/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package fs

import (
	"testing"
)

func TestFind(t *testing.T) {
	files, _ := Find("./*.go")

	if len(files) == 0 {
		t.Fatal()
	}
}

func TestGlob(t *testing.T) {
	if GlobMatch("*.py", "a.rb") {
		t.Fatal()
	}
	if !GlobMatch("*.rb", "a.rb") {
		t.Fatal()
	}
	if !GlobMatch(".", "a.rb") {
		t.Fatal()
	}
	if !GlobMatch("/**/*.rb", "/full/path/a.rb") {
		t.Fatal()
	}
}

func TestIs(t *testing.T) {
	if !IsFile("fs.go") {
		t.Fatal()
	}
	if IsDir("fs.go") {
		t.Fatal()
	}
}
