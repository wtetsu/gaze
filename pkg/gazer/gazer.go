/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

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
func (g *Gazer) Run(configs *config.Config, timeout int64, restart bool) error {
	if timeout <= 0 {
		return errors.New("timeout must be more than 0")
	}
	err := g.repeatRunAndWait(configs, timeout, restart)
	return err
}

func (g *Gazer) repeatRunAndWait(commandConfigs *config.Config, timeout int64, restart bool) error {
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
			rawCommandString, err := getAppropriateCommand(event.Name, commandConfigs)
			if err != nil {
				logger.NoticeObject(err)
				continue
			}

			queueManageKey := rawCommandString
			commandStringList := splitCommand(queueManageKey)
			if len(commandStringList) == 0 {
				logger.Debug("Command not found: %s", event.Name)
				continue
			}

			ongoingCommand := g.commands.get(queueManageKey)
			if ongoingCommand != nil {
				if restart {
					kill(ongoingCommand.cmd, "Restart")
					g.commands.update(queueManageKey, nil)
				} else {
					g.commands.enqueue(queueManageKey, event)
					continue
				}
			}

			go func() {
				lastLaunched := time.Now()

				commandSize := len(commandStringList)

				timeoutCh := time.After(timeout)
				for i, commandString := range commandStringList {
					if commandSize == 1 {
						logger.NoticeWithBlank("[%s]", commandString)
					} else {
						logger.NoticeWithBlank("[%s](%d/%d)", commandString, i+1, commandSize)
					}

					err := g.invokeOneCommand(commandString, queueManageKey, timeoutCh)
					if err != nil {
						if len(err.Error()) > 0 {
							logger.NoticeObject(err)
						}
						break
					}
				}
				// Handle waiting events
				queuedEvent := g.commands.dequeue(queueManageKey)
				if queuedEvent == nil {
					g.commands.update(queueManageKey, nil)
				} else {
					canAbolish := lastLaunched > queuedEvent.Time
					if canAbolish {
						logger.Debug("Abolish:%d, %d", lastLaunched, queuedEvent.Time)
					} else {
						// Requeue
						g.commands.update(queueManageKey, nil)
						g.notify.Requeue(*queuedEvent)
					}
				}
			}()
		case <-sigInt:
			isDisposed = true
			return nil
		}
	}
	return nil
}

func (g *Gazer) invokeOneCommand(commandString string, queueManageKey string, timeoutCh <-chan struct{}) error {
	cmd := createCommand(commandString)
	g.commands.update(queueManageKey, cmd)
	err := executeCommandOrTimeout(cmd, timeoutCh)
	return err
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

var newLines = regexp.MustCompile("\r\n|\n\r|\n|\r")

func splitCommand(commandString string) []string {
	var commandList []string
	for _, rawCmd := range newLines.Split(commandString, -1) {
		cmd := strings.TrimSpace(rawCmd)
		if len(cmd) > 0 {
			commandList = append(commandList, cmd)
		}
	}
	return commandList
}

// Counter returns the current execution counter
func (g *Gazer) Counter() uint64 {
	return g.counter
}
