/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"
	"path/filepath"

	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/notify"
	"github.com/wtetsu/gaze/pkg/time"
)

// Gazer gazes filesystem.
type Gazer struct {
	patterns []string
	notify   *notify.Notify
	isClosed bool
	counter  uint64
	commands commands
}

// New returns a new Gazer.
func New(patterns []string) *Gazer {
	cleanPatterns := make([]string, len(patterns))
	for i, p := range patterns {
		cleanPatterns[i] = filepath.Clean(p)
	}

	notify, _ := notify.New(cleanPatterns)
	return &Gazer{
		patterns: cleanPatterns,
		notify:   notify,
		isClosed: false,
		counter:  0,
		commands: newCommands(),
	}
}

// Close disposes internal resources.
func (g *Gazer) Close() {
	if g.isClosed {
		return
	}
	g.notify.Close()
	g.isClosed = true
}

// Run starts to gaze.
func (g *Gazer) Run(configs *config.Config, timeout int, restart bool) error {
	err := g.repeatRunAndWait(configs, timeout, restart)
	return err
}

func (g *Gazer) repeatRunAndWait(commandConfigs *config.Config, timeout int, restart bool) error {
	sigInt := sigIntChannel()

	isDisposed := false
	for {
		if isDisposed {
			break
		}
		select {
		case event := <-g.notify.Events:
			logger.Debug("Receive: %s", event.Name)
			if !matchAny(g.patterns, event.Name) {
				continue
			}

			g.counter++
			commandString, err := getAppropriateCommand(event.Name, commandConfigs)
			if err != nil {
				logger.NoticeObject(err)
				continue
			}
			if commandString == "" {
				logger.Debug("Command not found: %s", event.Name)
				continue
			}

			ongoingCommand := g.commands.get(commandString)
			if ongoingCommand != nil && !hasProcessExited(ongoingCommand.cmd) {
				if restart {
					kill(ongoingCommand.cmd, "Restart")
					g.commands.update(commandString, nil)
				} else {
					g.notify.Enqueue(commandString, event)
					continue
				}
			}
			logger.NoticeWithBlank("[%s]", commandString)

			cmd := createCommand(commandString)
			g.commands.update(commandString, cmd)
			go func() {
				lastLaunched := time.Now()
				err := executeCommandOrTimeout(cmd, timeout)
				if err != nil {
					logger.NoticeObject(err)
				}
				// Handle waiting events
				for {
					queuedEvent := g.notify.Dequeue(commandString)
					if queuedEvent == nil {
						break
					}
					canAbolish := lastLaunched > queuedEvent.Time
					if canAbolish {
						logger.Notice("Abolish:%d, %d", lastLaunched, queuedEvent.Time)
						continue
					}
					// Requeue
					g.notify.Requeue(*queuedEvent)
				}
			}()
		case <-sigInt:
			isDisposed = true
			return nil
		}
	}
	return nil
}

func hasProcessExited(cmd *exec.Cmd) bool {
	if cmd.ProcessState == nil {
		return false
	}
	return cmd.ProcessState.Exited()
}

func matchAny(watchFiles []string, s string) bool {
	result := false
	for _, f := range watchFiles {
		if fs.GlobMatch(f, s) {
			result = true
			break
		}
	}
	return result
}

func getAppropriateCommand(filePath string, commandConfigs *config.Config) (string, error) {
	var result string
	var resultError error
	for _, c := range commandConfigs.Commands {
		if c.Cmd == "" || c.Ext == "" && c.Re == "" {
			continue
		}
		if c.Match(filePath) {
			command, err := render(c.Cmd, filePath)
			result = command
			resultError = err
			break
		}
	}
	return result, resultError
}

// Counter returns the current execution counter
func (g *Gazer) Counter() uint64 {
	return g.counter
}
