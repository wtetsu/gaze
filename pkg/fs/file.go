package fs

import "os"

// Exist returns true if filePath exists.
func Exist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
