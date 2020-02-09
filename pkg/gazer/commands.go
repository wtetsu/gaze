/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"sync"

	"github.com/wtetsu/gaze/pkg/time"
)

type commands struct {
	commands     map[string]command
	commandMutex sync.Mutex
}

type command struct {
	cmd          *exec.Cmd
	lastLaunched int64
}

func newCommands() commands {
	return commands{
		commands: make(map[string]command),
	}
}

func (c *commands) update(key string, cmd *exec.Cmd) {
	c.commandMutex.Lock()
	defer c.commandMutex.Unlock()

	if cmd == nil {
		delete(c.commands, key)
		return
	}
	c.commands[key] = command{cmd: cmd, lastLaunched: time.Now()}
}

func (c *commands) get(key string) *command {
	c.commandMutex.Lock()
	defer c.commandMutex.Unlock()

	cmd, ok := c.commands[key]
	if !ok {
		return nil
	}
	return &cmd
}
