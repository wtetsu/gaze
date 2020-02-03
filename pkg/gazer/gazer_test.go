/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/time"
)

func TestBasic(t *testing.T) {
	py1 := createTempFile("*.py", `import datetime; print(datetime.datetime.now())`)
	rb1 := createTempFile("*.rb", `print(Time.new)`)

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer := New([]string{py1, rb1})
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.InitConfig([]string{".gaze.yml", ".gaze.yaml"})
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 0, false)
	if gazer.Counter() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(rb1)
		if gazer.Counter() >= 4 {
			break
		}
		time.Sleep(50)
	}

	if gazer.Counter() < 4 {
		t.Fatal()
	}
}

func TestRestart(t *testing.T) {
	content := `
import time

print("start")
# time.sleep(1)
print("end")
`

	py1 := createTempFile("*.py", content)
	if py1 == "" {
		t.Fatal("Temp files error")
	}

	gazer := New([]string{py1})
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.InitConfig([]string{".gaze.yml", ".gaze.yaml"})
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 0, true)

	if gazer.Counter() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(py1)
		touch(py1)
		if gazer.Counter() >= 2 {
			break
		}
		time.Sleep(10)
	}

	if gazer.Counter() < 2 {
		t.Fatalf("count:%d", gazer.Counter())
	}

	gazer.Close()
	gazer.Close()
}

func TestKill(t *testing.T) {
	py1 := createTempFile("*.py", `import time; time.sleep(5)`)
	rb1 := createTempFile("*.rb", `sleep(5)`)

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer := New([]string{py1, rb1})
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.InitConfig([]string{".gaze.yml", ".gaze.yaml"})
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 0, false)
	if gazer.Counter() != 0 {
		t.Fatal()
	}

	touch(py1)
	touch(rb1)

	pyKilled := false
	rbKilled := false
	for i := 0; i < 100; i++ {
		if !pyKilled && kill(gazer.commands[py1], "test") {
			pyKilled = true
		}
		if !rbKilled && kill(gazer.commands[rb1], "test") {
			rbKilled = true
		}
		if pyKilled && rbKilled {
			break
		}
		time.Sleep(10)
	}

	if !pyKilled || !rbKilled {
		t.Fatal()
	}
}

func TestGetAppropriateCommand(t *testing.T) {
	var commandConfigs config.Config

	var command string
	var err error

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: "", Cmd: "echo"})
	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".txt", Cmd: ""})

	command, err = getAppropriateCommand("a.txt", &commandConfigs)
	if command != "" {
		t.Fatal()
	}

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".txt", Cmd: "echo"})

	command, err = getAppropriateCommand("", &commandConfigs)
	if command == "a.txt" || err != nil {
		t.Fatal()
	}
}

func createTempFile(pattern string, content string) string {
	dirpath, err := ioutil.TempDir("", "_gaze")

	if err != nil {
		return ""
	}

	file, err := ioutil.TempFile(dirpath, pattern)
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
	file.WriteString("#\n")
	file.Close()
}
