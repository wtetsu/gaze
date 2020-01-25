/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package notify

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Notify delives events to a channel when files are virtually updated.
// "create+rename" is regarded as "update".
type Notify struct {
	Events        chan Event
	Errors        chan error
	watcher       *fsnotify.Watcher
	isClosed      bool
	times         map[string]int64
	pendingPeriod int64
}

// Event represents a single file system notification.
type Event struct {
	Name string
}

// Op describes a set of file operations.
type Op = fsnotify.Op

// Close disposes internal resources.
func (n *Notify) Close() {
	if n.isClosed {
		return
	}
	n.watcher.Close()
	n.isClosed = true
}

// New creates a Notify
func New(patterns []string) (*Notify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.ErrorObject(err)
		return nil, err
	}

	watchDirs := findDirs(patterns)

	for _, t := range watchDirs {
		err = watcher.Add(t)
		if err != nil {
			logger.Error("%s: %v", t, err)
		}
		logger.Info("gazing at: %s", t)
	}

	notify := &Notify{
		Events:        make(chan Event),
		watcher:       watcher,
		isClosed:      false,
		times:         make(map[string]int64),
		pendingPeriod: 100 * 1000000,
	}

	go notify.wait()

	return notify, nil
}

func findDirs(patterns []string) []string {
	targets := uniq.New()
	for _, pattern := range patterns {
		patternDir := filepath.Dir(pattern)
		if fs.IsDir(patternDir) {
			targets.Add(patternDir)
		}
		_, dirs := fs.Find(pattern)
		for _, d := range dirs {
			targets.Add(d)
		}
	}
	return targets.List()
}

func (n *Notify) wait() {
	for {
		select {
		case event, ok := <-n.watcher.Events:
			if !ok {
				continue
			}
			if !n.shouldExecute(event.Name, event.Op) {
				continue
			}
			n.times[event.Name] = time.Now()
			e := Event{
				Name: event.Name,
			}
			n.Events <- e
		case err, ok := <-n.watcher.Errors:
			if !ok {
				continue
			}
			n.Errors <- err
		}
	}
}

func (n *Notify) shouldExecute(filePath string, op Op) bool {
	flag := fsnotify.Write | fsnotify.Rename
	if op|flag == 0 {
		return false
	}

	lastExecutionTime := n.times[filePath]
	modifiedTime := time.GetFileModifiedTime(filePath)
	if (modifiedTime - lastExecutionTime) < n.pendingPeriod {
		return false
	}
	return true
}
