/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	"os"
	"path"
	"testing"

	"github.com/cbroglie/mustache" // 追加
)

func TestNew(t *testing.T) {
	config := New("ls")

	if !config.Commands[0].Match("abcdefg") {
		t.Fatal()
	}
}

func TestInit(t *testing.T) {
	InitConfig()
}

func TestMatch(t *testing.T) {
	yaml := createTempFile("*.yml", testConfig())
	c, err := makeConfigFromFile(yaml)

	if err != nil {
		t.Fatal(err)
	}
	if len(c.Commands) != 4 {
		t.Fatal()
	}
	if getFirstMatch(c, "") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "a.rb").Cmd != "run01" {
		t.Fatal()
	}
	if getFirstMatch(c, "Dockerfile").Cmd != "run02" {
		t.Fatal()
	}
	if getFirstMatch(c, ".Dockerfile") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "Dockerfile.") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "abc.txt").Cmd != "run03" {
		t.Fatal()
	}
	if getFirstMatch(c, "abcdef.txt").Cmd != "run03" {
		t.Fatal()
	}
	if getFirstMatch(c, "ab.txt") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "abc") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "zzz.txt") != nil {
		t.Fatal()
	}
}

func TestInvalidYaml(t *testing.T) {
	c, err := makeConfigFromBytes([]byte("aaa_bbb_ccc"))
	if err == nil {
		t.Fatal()
	}
	if c != nil {
		t.Fatal()
	}

	c, err = makeConfigFromFile("___.yml")
	t.Log(c)
	t.Log(err)
	if err == nil {
		t.Fatal()
	}
	if c != nil {
		t.Fatal()
	}
}

func TestSearchConfigPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "__gaze_test")
	if err != nil {
		t.Fatal(err)
	}

	if searchConfigPath("") != "" {
		t.Fatal()
	}

	// Should be not found
	if searchConfigPath(tempDir) != "" {
		t.Fatal()
	}

	os.Create(path.Join(tempDir, ".gaze.yaml"))
	if searchConfigPath(tempDir) != path.Join(tempDir, ".gaze.yaml") {
		t.Fatal()
	}

	os.Create(path.Join(tempDir, ".gaze.yml"))
	if searchConfigPath(tempDir) != path.Join(tempDir, ".gaze.yml") {
		t.Fatal()
	}

	os.MkdirAll(path.Join(tempDir, ".config", "gaze"), os.ModePerm)
	os.Create(path.Join(tempDir, ".config", "gaze", "gaze.yaml"))
	if searchConfigPath(tempDir) != path.Join(tempDir, ".config", "gaze", "gaze.yaml") {
		t.Fatal()
	}

	os.Create(path.Join(tempDir, ".config", "gaze", "gaze.yml"))
	if searchConfigPath(tempDir) != path.Join(tempDir, ".config", "gaze", "gaze.yml") {
		t.Fatal()
	}
}

func getFirstMatch(config *Config, fileName string) *Command {
	var result *Command
	for _, command := range config.Commands {
		if command.Match(fileName) {
			result = &command
			break
		}
	}
	return result
}

func testConfig() string {
	return `#
commands:
- ext:
  cmd: run00
- ext: .rb
  cmd: run01
- re: ^Dockerfile$
  cmd: run02
- re: ^abc
  ext: .txt
  cmd: run03
`
}

func createTempFile(pattern string, content string) string {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return ""
	}
	file.WriteString(content)
	file.Close()

	return file.Name()
}

func TestRenderStartEnd(t *testing.T) {
	logConf := &Log{
		Start: "Start: {{key}}",
		End:   "End: {{key}}",
	}
	startTmpl, err := mustache.ParseString(logConf.Start)
	if err != nil {
		t.Fatalf("failed to parse start template: %s", err)
	}
	logConf.start = startTmpl

	endTmpl, err := mustache.ParseString(logConf.End)
	if err != nil {
		t.Fatalf("failed to parse end template: %s", err)
	}
	logConf.end = endTmpl

	params := map[string]string{"key": "value1"}
	startResult := logConf.RenderStart(params)
	expectedStart := "Start: value1"
	if startResult != expectedStart {
		t.Fatalf("expected %q but got %q", expectedStart, startResult)
	}

	params = map[string]string{"key": "value2"}
	endResult := logConf.RenderEnd(params)
	expectedEnd := "End: value2"
	if endResult != expectedEnd {
		t.Fatalf("expected %q but got %q", expectedEnd, endResult)
	}

	startResult = logConf.RenderStart(nil)
	expectedStart = "Start: "
	if startResult != expectedStart {
		t.Fatalf("expected %q but got %q", expectedStart, startResult)
	}
}

func TestRenderLog(t *testing.T) {
	templateStr := "Hello, {{name}}!"
	tmpl, err := mustache.ParseString(templateStr)
	if err != nil {
		t.Fatalf("failed to parse template: %s", err)
	}
	params := map[string]string{
		"name": "World",
	}
	result := renderLog(tmpl, params)
	expected := "Hello, World!"
	if result != expected {
		t.Fatalf("expected %q but got %q", expected, result)
	}

	result = renderLog(tmpl, nil)
	expected = "Hello, !"
	if result != expected {
		t.Fatalf("expected %q but got %q", expected, result)
	}
}
