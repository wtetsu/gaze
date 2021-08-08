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
	"sync"

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
	mutexes  map[string]*sync.Mutex
}

// New returns a new Gazer.
func New(patterns []string, maxWatchDirs int) (*Gazer, error) {
	cleanPatterns := make([]string, len(patterns))
	for i, p := range patterns {
		cleanPatterns[i] = filepath.Clean(p)
	}

	notify, err := notify.New(cleanPatterns, maxWatchDirs)
	if err != nil {
		return nil, err
	}
	return &Gazer{
		patterns: cleanPatterns,
		notify:   notify,
		isClosed: false,
		counter:  0,
		commands: newCommands(),
		mutexes:  make(map[string]*sync.Mutex),
	}, nil
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
		select {
		case event := <-g.notify.Events:
			if isDisposed {
				break
			}
			logger.Debug("Receive: %s", event.Name)
			g.counter++

			commandStringList := g.tryToFindCommand(event.Name, commandConfigs)
			if commandStringList == nil {
				continue
			}

			queueManageKey := strings.Join(commandStringList, "\n")

			ongoingCommand := g.commands.get(queueManageKey)

			if ongoingCommand != nil && restart {
				kill(ongoingCommand.cmd, "Restart")
				g.commands.update(queueManageKey, nil)
			}

			if ongoingCommand != nil && !restart {
				g.commands.enqueue(queueManageKey, event)
				continue
			}

			mutex := g.lock(queueManageKey)

			go func() {
				g.invoke(commandStringList, queueManageKey, timeout)
				mutex.Unlock()
			}()

		case <-sigInt:
			isDisposed = true
			return nil
		}
	}
}

func (g *Gazer) tryToFindCommand(filePath string, commandConfigs *config.Config) []string {
	if !matchAny(g.patterns, filePath) {
		return nil
	}

	rawCommandString, err := getMatchedCommand(filePath, commandConfigs)
	if err != nil {
		logger.NoticeObject(err)
		return nil
	}

	commandStringList := splitCommand(rawCommandString)
	if len(commandStringList) == 0 {
		logger.Debug("Command not found: %s", filePath)
		return nil
	}

	return commandStringList
}

func (g *Gazer) lock(queueManageKey string) *sync.Mutex {
	mutex, ok := g.mutexes[queueManageKey]
	if !ok {
		mutex = &sync.Mutex{}
		g.mutexes[queueManageKey] = mutex
	}
	mutex.Lock()
	return mutex
}

func (g *Gazer) invoke(commandStringList []string, queueManageKey string, timeout int64) {
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
}

func (g *Gazer) invokeOneCommand(commandString string, queueManageKey string, timeoutCh <-chan struct{}) error {
	cmd := createCommand(commandString)
	g.commands.update(queueManageKey, cmd)
	err := executeCommandOrTimeout(cmd, timeoutCh)
	return err
}

func matchAny(watchFiles []string, s string) bool {
	for _, f := range watchFiles {
		if fs.GlobMatch(f, s) {
			return true
		}
	}
	return false
}

func getMatchedCommand(filePath string, commandConfigs *config.Config) (string, error) {
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
