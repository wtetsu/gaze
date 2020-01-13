/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	"testing"
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

func TestXxx(t *testing.T) {
	invalid, err := makeConfig([]byte(testConfig()))
	if err != nil {
		t.Fatal(err)
	}
	if len(invalid.Commands) != 4 {
		t.Fatal()
	}
	if getFirstMatch(invalid, "") != nil {
		t.Fatal()
	}
	if getFirstMatch(invalid, "a.rb").Run != "run01" {
		t.Fatal()
	}
	if getFirstMatch(invalid, "Dockerfile").Run != "run02" {
		t.Fatal()
	}
	if getFirstMatch(invalid, ".Dockerfile") != nil {
		t.Fatal()
	}
	if getFirstMatch(invalid, "Dockerfile.") != nil {
		t.Fatal()
	}
	if getFirstMatch(invalid, "abc.txt").Run != "run03" {
		t.Fatal()
	}
	if getFirstMatch(invalid, "abcdef.txt").Run != "run03" {
		t.Fatal()
	}
	if getFirstMatch(invalid, "ab.txt") != nil {
		t.Fatal()
	}
	if getFirstMatch(invalid, "abc") != nil {
		t.Fatal()
	}
	if getFirstMatch(invalid, "zzz.txt") != nil {
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

// Default returns the default configuration
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
