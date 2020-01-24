/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	"io/ioutil"
	"testing"
)

func TestNew(t *testing.T) {
	config := New("ls")

	if !config.Commands[0].Match("abcdefg") {
		t.Fatal()
	}
}

func TestInit(t *testing.T) {
	InitConfig([]string{".gaze.yml", ".gaze.yaml"})
}

func TestMatch(t *testing.T) {
	yaml := createTempFile("*.yml", testConfig())
	c, err := makeConfig(yaml)

	if err != nil {
		t.Fatal(err)
	}
	if len(c.Commands) != 4 {
		t.Fatal()
	}
	if getFirstMatch(c, "") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "a.rb").Run != "run01" {
		t.Fatal()
	}
	if getFirstMatch(c, "Dockerfile").Run != "run02" {
		t.Fatal()
	}
	if getFirstMatch(c, ".Dockerfile") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "Dockerfile.") != nil {
		t.Fatal()
	}
	if getFirstMatch(c, "abc.txt").Run != "run03" {
		t.Fatal()
	}
	if getFirstMatch(c, "abcdef.txt").Run != "run03" {
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
	c, err := makeConfig("___.yml")
	t.Log(c)
	t.Log(err)
	if err == nil {
		t.Fatal()
	}
	if c != nil {
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
  run: run00
- ext: .rb
  run: run01
- re: ^Dockerfile$
  run: run02
- re: ^abc
  ext: .txt
  run: run03
`
}

func createTempFile(pattern string, content string) string {
	file, err := ioutil.TempFile("", pattern)
	if err != nil {
		return ""
	}
	file.WriteString(content)
	file.Close()

	return file.Name()
}
