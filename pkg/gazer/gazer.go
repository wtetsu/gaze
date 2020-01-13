/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

// Gazer gazes filesystem.
type Gazer struct {
	patterns []string
	watcher  *fsnotify.Watcher
	isClosed bool
	counter  uint64
}

// New returns a new Gazer.
func New(patterns []string) *Gazer {
	watcher, _ := createWatcher(patterns)
	return &Gazer{
		patterns: patterns,
		watcher:  watcher,
		isClosed: false,
	}
}

// Close disposes internal resources.
func (g *Gazer) Close() {
	if g.isClosed {
		return
	}
	g.watcher.Close()
	g.isClosed = true
}

// Run starts to gaze.
func (g *Gazer) Run(configs *config.Config, timeout int, restart bool) error {
	err := g.repeatRunAndWait(configs, timeout, restart)
	return err
}

func createWatcher(patterns []string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.ErrorObject(err)
		return nil, err
	}

	added := map[string]struct{}{}
	for _, pattern := range patterns {
		_, dirs := fs.Find(pattern)
		for _, d := range dirs {
			_, ok := added[d]
			if ok {
				continue
			}
			logger.Info("gazing at: %s", d)
			err = watcher.Add(d)
			if err != nil {
				logger.ErrorObject(err)
			}
			added[d] = struct{}{}
		}
	}

	return watcher, nil
}

func (g *Gazer) repeatRunAndWait(commandConfigs *config.Config, timeout int, restart bool) error {
	var lastExecutionTime int64

	sigInt := sigIntChannel()

	var ignorePeriod int64 = 10 * 1000000

	var ongoingCommand *exec.Cmd

	isDisposed := false
	for {
		if isDisposed {
			break
		}
		select {
		case event, ok := <-g.watcher.Events:
			flag := fsnotify.Write | fsnotify.Rename
			if ok && event.Op|flag == 0 {
				continue
			}
			if !matchAny(g.patterns, event.Name) {
				continue
			}
			logger.Debug("Receive: %s", event.Name)
			modifiedTime := time.GetFileModifiedTime(event.Name)
			if (modifiedTime - lastExecutionTime) < ignorePeriod {
				continue
			}

			g.counter++
			commandString := getAppropriateCommand(event.Name, commandConfigs)
			if commandString != "" {
				logger.NoticeWithBlank("[%s]", commandString)

				if ongoingCommand != nil {
					kill(ongoingCommand, "Restart")
					ongoingCommand = nil
				}

				cmd := createCommand(commandString)
				lastExecutionTime = time.Now()
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
