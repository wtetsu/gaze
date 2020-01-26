/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
)

var templateCache = make(map[string]*mustache.Template)

func render(sourceString string, rawfilePath string) (string, error) {
	template, err := getOrCreateTemplate(sourceString)
	if err != nil {
		return "", err
	}

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

	return result, err
}

func getOrCreateTemplate(sourceString string) (*mustache.Template, error) {
	cachedTemplate, ok := templateCache[sourceString]
	if ok {
		return cachedTemplate, nil
	}

	template, err := mustache.ParseString(sourceString)
	if err != nil {
		return nil, fmt.Errorf("%v(%s)", err, sourceString)
	}
	templateCache[sourceString] = template
	return template, nil
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
