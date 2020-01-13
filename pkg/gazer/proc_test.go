/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"testing"
	"time"
)

func TestProc1(t *testing.T) {
	// Very normal
	cmd := createCommand("echo hello")
	executeCommandOrTimeout(cmd, 0)
}

func TestProc2(t *testing.T) {
	// Kill using timeout
	cmd := createCommand("sleep 60")
	executeCommandOrTimeout(cmd, 100)
}

func TestProc3(t *testing.T) {
	// Kill using a signal
	cmd := createCommand("sleep 60")
	go executeCommandOrTimeout(cmd, 0)

	for {
		time.Sleep(50)
		if cmd.Process != nil {
			cmd.Process.Kill()
			break
		}
	}
}
