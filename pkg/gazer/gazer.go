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
)

// Gazer gazes filesystem.
type Gazer struct {
	patterns []string
	notify   *notify.Notify
	isClosed bool
	counter  uint64
	commands map[string]*exec.Cmd
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
		commands: make(map[string]*exec.Cmd),
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
	var ongoingCommand *exec.Cmd

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
			commandString := getAppropriateCommand(event.Name, commandConfigs)
			if commandString == "" {
				logger.Debug("Command not found: %s", event.Name)
				continue
			}

			logger.NoticeWithBlank("[%s]", commandString)

			if ongoingCommand != nil {
				kill(ongoingCommand, "Restart")
				ongoingCommand = nil
			}

			cmd := createCommand(commandString)
			g.commands[event.Name] = cmd
			if !restart {
				err := executeCommandOrTimeout(cmd, timeout)
				if err != nil {
					logger.NoticeObject(err)
				}
			} else {
				// restartable
				ongoingCommand = cmd
				go func() {
					err := executeCommandOrTimeout(cmd, timeout)
					if err != nil {
						logger.NoticeObject(err)
					}
					ongoingCommand = nil
				}()
			}

		case <-sigInt:
			isDisposed = true
			return nil
		}
	}
	return nil
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

func getAppropriateCommand(filePath string, commandConfigs *config.Config) string {
	var result string
	for _, c := range commandConfigs.Commands {
		if c.Run == "" || c.Ext == "" && c.Re == "" {
			continue
		}
		if c.Match(filePath) {
			command := render(c.Run, filePath)
			result = command
			break
		}
	}
	return result
}

// Counter returns the current execution counter
func (g *Gazer) Counter() uint64 {
	return g.counter
}
