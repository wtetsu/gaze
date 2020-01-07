package file

import "os"

// Exist returns true if filePath exists.
func Exist(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
