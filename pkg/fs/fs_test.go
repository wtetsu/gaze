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

	if !GlobMatch("xx9", "xx9") {
		t.Fatal()
	}
	if !GlobMatch("xx9", "xx9/a.rb") {
		t.Fatal()
	}

	if !GlobMatch("xx?yy", "xx9yy") {
		t.Fatal()
	}
	if !GlobMatch("xx?yy", "xx9yy/a.rb") {
		t.Fatal()
	}
	if GlobMatch("xx?yy", "xx99yy") {
		t.Fatal()
	}
	if GlobMatch("xx?yy", "xx99yy/a.rb") {
		t.Fatal()
	}

	if !GlobMatch("xx*yy", "xx9yy") {
		t.Fatal()
	}
	if !GlobMatch("xx*yy", "xx9yy/a.rb") {
		t.Fatal()
	}
	if !GlobMatch("xx*yy", "xx99yy") {
		t.Fatal()
	}
	if !GlobMatch("xx*yy", "xx99yy/a.rb") {
		t.Fatal()
	}
}

func TestIs(t *testing.T) {
	if !IsFile("fs.go") {
		t.Fatal()
	}
	if IsFile("__fs.go") {
		t.Fatal()
	}
	if !IsDir(".") || !IsDir("..") {
		t.Fatal()
	}
	if IsDir("fs.go") {
		t.Fatal()
	}
	if IsDir("__fs.go") {
		t.Fatal()
	}
}
