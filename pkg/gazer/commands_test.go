/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"testing"

	"github.com/wtetsu/gaze/pkg/notify"
	"github.com/wtetsu/gaze/pkg/tutil"
)

func TestCommandsBasic1(t *testing.T) {
	commands := newCommands()

	key := "key01"

	if commands.get(key) != nil {
		t.Fatal()
	}
	var cmd exec.Cmd
	commands.update(key, &cmd)
	if commands.get(key) == nil {
		t.Fatal()
	}
	commands.update(key, nil)
	if commands.get(key) != nil {
		t.Fatal()
	}
}

func TestCommandsBasic2(t *testing.T) {
	rbCommand := `ruby a.rb`
	pyCommand := `python a.py`

	commands := newCommands()

	if commands.dequeue(rbCommand) != nil {
		t.Fatal()
	}
	if commands.dequeue(pyCommand) != nil {
		t.Fatal()
	}
	commands.enqueue(rbCommand, notify.Event{Name: rbCommand, Time: 1})
	commands.enqueue(pyCommand, notify.Event{Name: pyCommand, Time: 2})

	if commands.dequeue(rbCommand) == nil {
		t.Fatal()
	}
	if commands.dequeue(pyCommand) == nil {
		t.Fatal()
	}
	if commands.dequeue(rbCommand) != nil {
		t.Fatal()
	}
	if commands.dequeue(pyCommand) != nil {
		t.Fatal()
	}
}

func TestCommandsParallel(t *testing.T) {
	commands := newCommands()

	key := "key01"

	go func() {
		for i := 0; i < 100; i++ {
			commands.get(key)
			tutil.Sleep(1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			var cmd exec.Cmd
			commands.update(key, &cmd)
			tutil.Sleep(1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			commands.update(key, nil)
			tutil.Sleep(1)
		}
	}()

}
