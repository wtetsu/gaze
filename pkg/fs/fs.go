/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package fs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Find returns a list of files and directories that match the pattern.
func Find(pattern string) ([]string, []string) {
	foundFiles, err := doublestar.Glob(pattern)
	if err != nil {
		return []string{}, []string{}
	}
	entryList := append([]string{pattern}, foundFiles...)

	fileList, dirList := doFileDir(entryList)
	return fileList, dirList
}

func doFileDir(entries []string) ([]string, []string) {
	fileUniq := uniq.New()
	dirUniq := uniq.New()

	for _, entry := range entries {
		stat := Stat(entry)
		if stat == nil {
			continue
		}
		if isDir(entry) {
			dirUniq.Add(filepath.Clean(entry))
		} else {
			fileUniq.Add(entry)
			dirPath := filepath.Dir(entry)
			dirUniq.Add(dirPath)
		}
	}
	return fileUniq.List(), dirUniq.List()
}

func isDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// GlobMatch returns true if a pattern matches a path string
func GlobMatch(rawPattern string, rawFilePath string) bool {
	pattern := filepath.ToSlash(rawPattern)
	filePath := strings.TrimSuffix(filepath.ToSlash(rawFilePath), "/")

	ok, _ := doublestar.Match(pattern, filePath)
	if ok {
		logger.Debug("rawPattern:%s, rawFilePath:%s, true(file)", rawPattern, rawFilePath)
		return true
	}

	dirPath := filepath.ToSlash(filepath.Dir(filePath))

	ok, _ = doublestar.Match(pattern, dirPath)
	if ok {
		logger.Debug("rawPattern:%s, rawFilePath:%s, true(dir)", rawPattern, rawFilePath)
		return true
	}

	logger.Debug("rawPattern:%s, rawFilePath:%s, false", rawPattern, rawFilePath)
	return false
}

// IsDir returns true if path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile returns true if path is a file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Stat returns a FileInfo.
func Stat(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return info
}
