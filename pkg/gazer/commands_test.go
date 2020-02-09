/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"testing"

	"github.com/wtetsu/gaze/pkg/time"
)

func TestCommandsBasic(t *testing.T) {
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

func TestCommandsParallel(t *testing.T) {
	commands := newCommands()

	key := "key01"

	go func() {
		for i := 0; i < 100; i++ {
			commands.get(key)
			time.Sleep(1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			var cmd exec.Cmd
			commands.update(key, &cmd)
			time.Sleep(1)
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			commands.update(key, nil)
			time.Sleep(1)
		}
	}()

}
