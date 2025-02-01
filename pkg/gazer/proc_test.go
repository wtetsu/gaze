/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"math"
	"os/exec"
	"testing"

	"github.com/wtetsu/gaze/pkg/time"
)

func TestProc1(t *testing.T) {
	// Very normal
	cmd := createCommand("echo hello")
	executeCommandOrTimeout(cmd, time.After(math.MaxInt64))
}

func TestProc2(t *testing.T) {
	// Kill using timeout
	cmd := createCommand("sleep 60")
	executeCommandOrTimeout(cmd, time.After(100))
}

func TestProc3(t *testing.T) {
	// Kill using a signal
	cmd := createCommand("sleep 60")
	go executeCommandOrTimeout(cmd, time.After(60*10000))

	for {
		time.Sleep(50)
		if cmd.Process != nil {
			cmd.Process.Kill()
			break
		}
	}
}

func TestProc4(t *testing.T) {
	cmd1 := createCommand("ls")
	if len(cmd1.Args) != 1 {
		t.Fatal()
	}
	if kill(cmd1, "test") {
		t.Fatal()
	}

	cmd2 := createCommand("ls aaa.txt")
	if len(cmd2.Args) != 2 {
		t.Fatal()
	}
	if kill(cmd2, "test") {
		t.Fatal()
	}

	cmd3 := createCommand(`ls aaa.txt "Program Files"`)
	if len(cmd3.Args) != 3 {
		t.Fatal()
	}
	if kill(cmd3, "test") {
		t.Fatal()
	}
}

func TestProc5(t *testing.T) {
	var cmd *exec.Cmd = nil
	_, err := executeCommandOrTimeout(cmd, time.After(100))
	if err == nil {
		t.Fatal()
	}
}
