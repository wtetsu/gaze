package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
)

func ListDir(root string, pattern string) []string {
	result := []string{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// return err
				// log.Println(err)
				return nil
			}
			// fmt.Println(path)
			if !isDir(path) {
				return nil
			}
			if globMatch(pattern, path) {
				// fmt.Println(path, info.Size())
				result = append(result, path)
			} else {
				fmt.Println("    " + path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return result
}

func isDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func globMatch(pattern string, rawPath string) bool {
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
