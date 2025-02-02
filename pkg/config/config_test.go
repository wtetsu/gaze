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

func TestNewWithFixedCommand(t *testing.T) {
	config, _ := NewWithFixedCommand("ls")
	if !config.Commands[0].Match("abcdefg") {
		t.Fatal()
	}

	config, err := NewWithFixedCommand("")
	if err == nil {
		t.Fatal()
	}
}

func TestInit(t *testing.T) {
	LoadPreferredConfig()
}

func TestMatch(t *testing.T) {
	yaml := createTempFile("*.yml", testConfig())
	c, err := LoadConfigFromFile(yaml)

	if err != nil {
		t.Fatal(err)
	}
	if len(c.Commands) != 3 {
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
	rawConfig, err := parseRawConfigFromBytes([]byte("aaa_bbb_ccc"))
	if err == nil {
		t.Fatal()
	}
	if rawConfig != nil {
		t.Fatal()
	}

	config, err := LoadConfigFromFile("___.yml")
	t.Log(config)
	t.Log(err)
	if err == nil {
		t.Fatal()
	}
	if config != nil {
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
	logConf := &Log{}

	startTmpl, err := mustache.ParseString("Start: {{key}}")
	if err != nil {
		t.Fatalf("failed to parse start template: %s", err)
	}
	logConf.start = startTmpl

	endTmpl, err := mustache.ParseString("End: {{key}}")
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

func TestToConfig(t *testing.T) {
	rawCfg := &rawConfig{
		Commands: []rawCommand{
			// 1. Valid command with ext only.
			{Ext: ".go", Cmd: "runGo"},
			// 2. Valid command with ext and valid regexp.
			{Ext: ".rb", Re: "^test", Cmd: "runRb"},
			// 3. Both ext and re empty; should be skipped.
			{Cmd: "badCmd"},
			// 4. Invalid regexp; should be skipped.
			{Re: "(", Cmd: "invalidRegex"},
			// 5. Valid command with regexp only.
			{Re: "^match", Cmd: "matchCmd"},
			// 6. ext provided as empty while re is empty; should be skipped.
			{Ext: "", Cmd: "extOnly"},
			// 7. Valid command with regexp only.
			{Re: "^noExt", Cmd: "regexOnly"},
			// 8. Empty command; should be skipped.
			{Cmd: ""},
		},
		Log: &rawLog{
			Start: "start: {{var}}",
			End:   "end: {{var}}",
		},
	}

	cfg := toConfig(rawCfg)
	// Expected commands: #1, #2, #5, and #7.
	expectedCmdCount := 4
	if len(cfg.Commands) != expectedCmdCount {
		t.Fatalf("expected %d commands but got %d", expectedCmdCount, len(cfg.Commands))
	}

	// Test first command.
	cmd := cfg.Commands[0]
	if cmd.Cmd != "runGo" || cmd.Ext != ".go" || cmd.re != nil {
		t.Errorf("unexpected command 0: %+v", cmd)
	}

	// Test second command.
	cmd = cfg.Commands[1]
	if cmd.Cmd != "runRb" || cmd.Ext != ".rb" || cmd.re == nil {
		t.Errorf("unexpected command 1: %+v", cmd)
	} else {
		// Verify the regexp compiles and matches an example string.
		if !cmd.re.MatchString("test_file.rb") {
			t.Errorf("regex did not match expected string for command 1")
		}
	}

	// Test third command.
	cmd = cfg.Commands[2]
	if cmd.Cmd != "matchCmd" || cmd.re == nil {
		t.Errorf("unexpected command 2: %+v", cmd)
	}

	// Test fourth command.
	cmd = cfg.Commands[3]
	if cmd.Cmd != "regexOnly" || cmd.re == nil {
		t.Errorf("unexpected command 3: %+v", cmd)
	}

	// Test log rendering.
	if cfg.Log == nil {
		t.Fatal("expected Log not to be nil")
	}
	startOut := cfg.Log.RenderStart(map[string]string{"var": "X"})
	if startOut != "start: X" {
		t.Errorf("expected start log 'start: X', got '%s'", startOut)
	}
	endOut := cfg.Log.RenderEnd(map[string]string{"var": "Y"})
	if endOut != "end: Y" {
		t.Errorf("expected end log 'end: Y', got '%s'", endOut)
	}
}
func TestParseRawConfigFromFile(t *testing.T) {
	// Test with valid YAML content
	validYAML := `
commands:
- ext: .go
  cmd: runGo
log:
  start: "start: {{var}}"
  end: "end: {{var}}"
`
	tmpFile, err := os.CreateTemp("", "valid-config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(validYAML); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	rawCfg, err := parseRawConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error parsing valid YAML: %v", err)
	}
	if rawCfg == nil {
		t.Fatal("expected non-nil rawConfig")
	}
	if len(rawCfg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(rawCfg.Commands))
	}
	if rawCfg.Log == nil {
		t.Fatal("expected non-nil Log in rawConfig")
	}

	// Test with non-existent file path
	_, err = parseRawConfigFromFile("nonexistentfile.yml")
	if err == nil {
		t.Fatal("expected error when file does not exist")
	}

	// Test with invalid YAML content
	invalidYAML := "invalid: : yaml: ::::"
	tmpFile2, err := os.CreateTemp("", "invalid-config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile2.Name())
	if _, err := tmpFile2.WriteString(invalidYAML); err != nil {
		t.Fatal(err)
	}
	tmpFile2.Close()

	rawCfg, err = parseRawConfigFromFile(tmpFile2.Name())
	if err == nil {
		t.Fatal("expected error when YAML is invalid")
	}
	if rawCfg != nil {
		t.Fatal("expected nil rawConfig when YAML is invalid")
	}
}

func TestLoadPreferredRawConfigFound(t *testing.T) {
	// Create a temporary directory to act as the home directory.
	tempHome, err := os.MkdirTemp("", "gaze-test-found")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	// Prepare a configuration file in the prioritized directory:
	// .config/gaze/gaze.yml
	configDir := path.Join(tempHome, ".config", "gaze")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	configFilePath := path.Join(configDir, "gaze.yml")
	configContent := `commands:
- ext: .test
  cmd: testCmd
log:
  start: "start: {{var}}"
  end: "end: {{var}}"
`
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	rawCfg, err := loadPreferredRawConfig(tempHome)
	if err != nil {
		t.Fatalf("loadPreferredRawConfig returned error: %s", err)
	}
	if rawCfg == nil {
		t.Fatal("expected non-nil rawConfig")
	}
	if len(rawCfg.Commands) == 0 {
		t.Fatal("expected at least one command from the configuration file")
	}
	if rawCfg.Commands[0].Cmd != "testCmd" {
		t.Fatalf("expected command 'testCmd', got %q", rawCfg.Commands[0].Cmd)
	}
}

func TestLoadPreferredRawConfigDefault(t *testing.T) {
	// Create a temporary directory with no config file.
	tempHome, err := os.MkdirTemp("", "gaze-test-default")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	// Ensure there is no config file anywhere under tempHome.
	// loadPreferredRawConfig will fall back to using the default configuration.
	rawCfg, err := loadPreferredRawConfig(tempHome)
	if err != nil {
		t.Fatalf("loadPreferredRawConfig returned error: %s", err)
	}
	if rawCfg == nil {
		t.Fatal("expected non-nil rawConfig")
	}
	// Since default configuration is used, we expect Log to be non-nil.
	if rawCfg.Log == nil {
		t.Fatal("expected non-nil Log in the default configuration")
	}
}
