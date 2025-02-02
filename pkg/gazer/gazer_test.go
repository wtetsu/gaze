/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/wtetsu/gaze/pkg/config"
)

func TestBasic(t *testing.T) {
	py1 := createTempFile("*.py", `import datetime; print(datetime.datetime.now())`)
	rb1 := createTempFile("*.rb", `print(Time.new)`)

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer, _ := New([]string{py1, rb1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.LoadPreferredConfig()
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 10*1000, false)
	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(rb1)
		if gazer.InvokeCount() >= 4 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if gazer.InvokeCount() < 4 {
		t.Fatal()
	}
}

func TestDoNothing(t *testing.T) {
	py1 := createTempFile("a'aa*.py", `import datetime; print(datetime.datetime.now())`)
	rb1 := createTempFile("b'bb.*.rb", `print(Time.new)`)

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer, _ := New([]string{py1, rb1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.LoadPreferredConfig()
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 10*1000, false)
	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(rb1)
		if gazer.InvokeCount() >= 4 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	if gazer.InvokeCount() > 0 {
		t.Fatal()
	}
}

func TestRename(t *testing.T) {
	py1 := createTempFile("*.py", `import datetime; print(datetime.datetime.now())`)
	rb1 := createTempFile("*.rb", `print(Time.new)`)

	py2 := py1 + ".tmp"
	rb2 := rb1 + ".tmp"

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer, _ := New([]string{py1, rb1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.LoadPreferredConfig()
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 10*1000, false)
	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	for i := 0; i < 20; i++ {
		if gazer.InvokeCount() >= 10 {
			break
		}

		os.Rename(py1, py2)
		os.Rename(rb1, rb2)

		time.Sleep(50 * time.Millisecond)

		touch(py1)
		os.Rename(py2, py1)
		touch(rb2)
		os.Rename(rb2, rb1)

		time.Sleep(50 * time.Millisecond)
	}

	if gazer.InvokeCount() < 10 {
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

	gazer, _ := New([]string{py1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.LoadPreferredConfig()
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 10*1000, true)

	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(py1)
		touch(py1)
		if gazer.InvokeCount() >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if gazer.InvokeCount() < 2 {
		t.Fatalf("count:%d", gazer.InvokeCount())
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

	gazer, _ := New([]string{py1, rb1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.LoadPreferredConfig()
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 10*1000, false)
	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	touch(py1)
	touch(rb1)

	py1Command := fmt.Sprintf(`python "%s"`, py1)
	rb1Command := fmt.Sprintf(`ruby "%s"`, rb1)

	pyKilled := false
	rbKilled := false
	for i := 0; i < 100; i++ {
		if !pyKilled && kill(getCmd(&gazer.commands, py1Command), "test") {
			pyKilled = true
		}
		if !rbKilled && kill(getCmd(&gazer.commands, rb1Command), "test") {
			rbKilled = true
		}
		if pyKilled && rbKilled {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !pyKilled || !rbKilled {
		t.Fatal()
	}
}

func getCmd(commands *commands, command string) *exec.Cmd {
	c := commands.get(command)
	if c == nil {
		return nil
	}

	return c.cmd
}

func TestInvalidCommand(t *testing.T) {
	py1 := createTempFile("*.py", `import datetime; print(datetime.datetime.now())`)
	rb1 := createTempFile("*.rb", `print(Time.new)`)

	if py1 == "" || rb1 == "" {
		t.Fatal("Temp files error")
	}

	gazer, _ := New([]string{py1, rb1}, 100)
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	var commandConfigs config.Config

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".rb", Cmd: "ruby {{file]]"})
	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".py", Cmd: ""})

	go gazer.Run(&commandConfigs, 10*1000, false)
	if gazer.InvokeCount() != 0 {
		t.Fatal()
	}

	for i := 0; i < 100; i++ {
		touch(py1)
		touch(rb1)
		if gazer.InvokeCount() >= 1 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	if gazer.InvokeCount() > 0 {
		t.Fatal()
	}
}

func TestGetAppropriateCommandOk(t *testing.T) {
	var commandConfigs config.Config

	var command string
	var err error

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: "", Cmd: "echo"})
	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".txt", Cmd: ""})

	command, err = getMatchedCommand("a.txt", commandConfigs.Commands)
	if command != "" || err != nil {
		t.Fatal()
	}

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".txt", Cmd: "echo"})

	command, err = getMatchedCommand("", commandConfigs.Commands)
	if command == "a.txt" || err != nil {
		t.Fatal()
	}
}

func TestGetAppropriateCommandError(t *testing.T) {
	var commandConfigs config.Config

	var command string
	var err error

	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".rb", Cmd: "ruby {{file]]"})
	commandConfigs.Commands = append(commandConfigs.Commands, config.Command{Ext: ".py", Cmd: "python {{file]]"})

	command, err = getMatchedCommand("a.txt", commandConfigs.Commands)
	if command != "" || err != nil {
		t.Fatal()
	}

	command, err = getMatchedCommand("a.rb", commandConfigs.Commands)
	if command != "" || err == nil {
		t.Fatal()
	}
	command, err = getMatchedCommand("a.py", commandConfigs.Commands)
	if command != "" || err == nil {
		t.Fatal()
	}
}

func TestInvalidTimeout(t *testing.T) {
	gazer, _ := New([]string{}, 100)

	var err error

	err = gazer.Run(nil, 0, false)
	if err == nil {
		t.Fatal()
	}

	err = gazer.Run(nil, -1, false)
	if err == nil {
		t.Fatal()
	}
}

func createTempFile(pattern string, content string) string {
	dirpath, err := os.MkdirTemp("", "_gaze")

	if err != nil {
		return ""
	}

	file, err := os.CreateTemp(dirpath, pattern)
	if err != nil {
		return ""
	}
	file.WriteString(content)
	file.Close()

	return filepath.ToSlash(file.Name())
}

func touch(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	file.WriteString("#\n")
	file.Close()
}
