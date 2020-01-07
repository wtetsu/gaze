package proc

import (
	"path/filepath"

	"github.com/cbroglie/mustache"
)

func render(sourceString string, filePath string) string {
	template, _ := mustache.ParseString(sourceString)

	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)
	abs, _ := filepath.Abs(filePath)
	dir := filepath.Dir(filePath)

	params := map[string]string{
		"file": filePath,
		"ext":  ext,
		"base": base,
		"abs":  abs,
		"dir":  dir,
	}

	result, err := template.Render(params)

	if err != nil {
		return ""
	}

	return result
}
