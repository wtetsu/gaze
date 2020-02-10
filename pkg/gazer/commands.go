/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"sync"

	"github.com/wtetsu/gaze/pkg/notify"
	"github.com/wtetsu/gaze/pkg/time"
)

type commands struct {
	commands map[string]command
	events   map[string]notify.Event
	mutex    sync.Mutex
}

type command struct {
	cmd          *exec.Cmd
	lastLaunched int64
	event        notify.Event
}

func newCommands() commands {
	return commands{
		commands: make(map[string]command),
		events:   make(map[string]notify.Event),
	}
}

func (c *commands) update(key string, cmd *exec.Cmd) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if cmd == nil {
		delete(c.commands, key)
		return
	}
	c.commands[key] = command{cmd: cmd, lastLaunched: time.Now()}
}

func (c *commands) get(key string) *command {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cmd, ok := c.commands[key]
	if !ok {
		return nil
	}
	return &cmd
}

func (c *commands) enqueue(commandString string, event notify.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.events[commandString] = event
}

func (c *commands) dequeue(commandString string) *notify.Event {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event, ok := c.events[commandString]

	if !ok {
		return nil
	}

	// delete both event and command
	delete(c.events, commandString)
	delete(c.commands, commandString)
	return &event
}
