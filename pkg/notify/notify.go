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
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Notify delives events to a channel when files are virtually updated.
// "create+rename" is regarded as "update".
type Notify struct {
	Events   chan Event
	Errors   chan error
	watcher  *fsnotify.Watcher
	isClosed bool
}

// Event represents a single file system notification.
type Event struct {
	Name string
	Op   Op
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
		Events:   make(chan Event),
		watcher:  watcher,
		isClosed: false,
	}

	go wait(notify)

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

func wait(notify *Notify) {
	for {
		select {
		case event, ok := <-notify.watcher.Events:
			if !ok {
				continue
			}
			e := Event{
				Name: event.Name,
				Op:   event.Op,
			}
			notify.Events <- e
		case err, ok := <-notify.watcher.Errors:
			if !ok {
				continue
			}
			notify.Errors <- err
		}
	}
}
