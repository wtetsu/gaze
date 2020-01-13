/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package app

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCreateCommandConfig(t *testing.T) {
	commandConfigs, err := createCommandConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if commandConfigs == nil || len(commandConfigs.Commands) == 0 {
		t.Fatal()
	}
}

func TestCreateCommandConfigWithUserCommand(t *testing.T) {
	commandConfigs, err := createCommandConfig("ls", "")
	if err != nil {
		t.Fatal(err)
	}
	if commandConfigs == nil || len(commandConfigs.Commands) != 1 {
		t.Fatal()
	}
}

func TestCreateCommandConfigWithFile(t *testing.T) {

	commandConfigs, err := createCommandConfig("", "no.yml")
	if commandConfigs != nil || err == nil {
		t.Fatal(err)
	}

	ymlFile := createTempFile("*.yml", yaml())

	commandConfigs, err = createCommandConfig("", ymlFile)
	if commandConfigs == nil || err != nil {
		t.Fatal(err)
	}
	if commandConfigs.Commands[0].Ext != ".py" {
		t.Fatal()
	}
	if commandConfigs.Commands[1].Ext != ".rb" {
		t.Fatal()
	}
	if commandConfigs.Commands[2].Ext != ".js" {
		t.Fatal()
	}
}

func TestEndTopEnd(t *testing.T) {
	rb := createTempFile("*.rb", `puts "Hello from Ruby`)
	py := createTempFile("*.py", `print("Hello from Python")`)

	watchFiles := []string{rb, py}
	userCommand := ""
	file := ""
	timeout := 0
	restart := false

	go Start(watchFiles, userCommand, file, timeout, restart)

	time.Sleep(100)
	touch(rb)
	time.Sleep(100)
	touch(py)
	time.Sleep(300)
}

func createTempFile(pattern string, content string) string {
	file, err := ioutil.TempFile("", "*.rb")
	if err != nil {
		return ""
	}
	file.WriteString(content)
	file.Close()

	return file.Name()
}

func touch(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	file.WriteString("")
	file.Close()
}

func yaml() string {
	return `#
commands:
- ext: .py
  run: python "{{file}}"
- ext: .rb
  run: ruby "{{file}}"
- ext: .js
  run: node "{{file}}"
`
}
