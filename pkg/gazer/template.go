/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"path/filepath"

	"github.com/cbroglie/mustache"
)

func render(sourceString string, rawfilePath string) string {
	template, _ := mustache.ParseString(sourceString)

	filePath := filepath.ToSlash(rawfilePath)
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)
	rawAbs, _ := filepath.Abs(filePath)
	abs := filepath.ToSlash(rawAbs)
	dir := filepath.ToSlash(filepath.Dir(filePath))

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
