/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gutil

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestGlobFuncWithError(t *testing.T) {
	// Define a custom glob function that returns an error.
	errorGlob := func(pattern string) ([]string, error) {
		return nil, fmt.Errorf("dummy error")
	}
	files, dirs := find("dummy", errorGlob)
	if len(files) != 0 || len(dirs) != 0 {
		t.Fatal("expected empty lists when glob function returns an error")
	}
}

func TestGlobFuncWithCustomSuccess(t *testing.T) {
	// Create a temporary file so that os.Stat can return valid info.
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	f.Close()

	// Define a custom glob function that returns our test file.
	customGlob := func(pattern string) ([]string, error) {
		return []string{filePath}, nil
	}

	files, dirs := find(filePath, customGlob)

	// Check that the file is listed.
	if len(files) == 0 {
		t.Fatal("expected test file in files list")
	}
	found := false
	for _, f := range files {
		if f == filePath {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("test file not found in files list")
	}

	// Check that the file's directory is included in the directory list.
	dirFound := false
	fileDir := filepath.Dir(filePath)
	for _, d := range dirs {
		if filepath.Clean(d) == filepath.Clean(fileDir) {
			dirFound = true
			break
		}
	}
	if !dirFound {
		t.Fatal("directory of test file not found in dirs list")
	}
}
func TestDoFileDir(t *testing.T) {
	// Create a temporary directory for testing.
	tmpDir := t.TempDir()

	// Create a temporary file inside tmpDir.
	filePath := filepath.Join(tmpDir, "test.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	f.Close()

	// Create a temporary subdirectory inside tmpDir.
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Build the entries list including:
	// - A file entry (twice to test uniqueness)
	// - A directory entry
	// - A non-existent file entry (which should be ignored)
	entries := []string{
		filePath,
		subDir,
		filePath,
		filepath.Join(tmpDir, "nonexistent"),
	}

	files, dirs := doFileDir(entries)

	// Verify that the file list contains exactly one element: filePath.
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0] != filePath {
		t.Fatalf("expected file path %s, got %s", filePath, files[0])
	}

	// The expected directories are:
	// - tmpDir (the parent of filePath)
	// - subDir (from the directory entry)
	expectedDirs := map[string]struct{}{
		filepath.Clean(tmpDir): {},
		filepath.Clean(subDir): {},
	}

	// Check that each expected directory is present in the results.
	for _, d := range dirs {
		cleaned := filepath.Clean(d)
		delete(expectedDirs, cleaned)
	}

	if len(expectedDirs) != 0 {
		t.Fatal("not all expected directories were found in the results")
	}
}
