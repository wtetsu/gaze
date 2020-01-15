/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"path/filepath"
	"strings"

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

	arr := strings.Split(base, ".")
	base0 := baseN(arr, 0)
	base1 := baseN(arr, 1)
	base2 := baseN(arr, 2)

	params := map[string]string{
		"file":  filePath,
		"ext":   ext,
		"base":  base,
		"abs":   abs,
		"dir":   dir,
		"base0": base0,
		"base1": base1,
		"base2": base2,
	}

	result, err := template.Render(params)

	if err != nil {
		return ""
	}

	return result
}

func baseN(arr []string, lastIndex int) string {
	var list []string
	for i := 0; ; i++ {
		if i > lastIndex || i >= len(arr) {
			break
		}
		list = append(list, arr[i])
	}
	return strings.Join(list, ".")
}
