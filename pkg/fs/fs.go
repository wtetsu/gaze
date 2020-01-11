package fs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Find returns
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
		if !Exist(entry) {
			continue
		}
		if isDir(entry) {
			dirUniq.Add(entry)
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
func GlobMatch(rawPattern string, rawPath string) bool {
	pattern := filepath.ToSlash(rawPattern)
	path := trimSuffix(filepath.ToSlash(rawPath), "/")

	ok, _ := doublestar.Match(pattern, path)
	if ok {
		return true
	}

	dirPath := filepath.Dir(pattern)
	ok, _ = doublestar.Match(dirPath, path)
	if ok {
		return true
	}
	return false
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func doublestarMatch(pattern string, path string) bool {
	ok, _ := doublestar.Match(pattern, path)
	if ok {
		return true
	}
	return false
}
