/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package notify

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Notify delivers events to a channel when files are virtually updated.
// "create+rename" is regarded as "update".
type Notify struct {
	Events                  chan Event
	Errors                  chan error
	watcher                 *fsnotify.Watcher
	isClosed                bool
	times                   map[string]int64
	pendingPeriod           int64
	regardRenameAsModPeriod int64
	detectCreate            bool
	candidates              []string
}

// Event represents a single file system notification.
type Event struct {
	Name string
	Time int64
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
func New(patterns []string, maxWatchDirs int) (*Notify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.ErrorObject(err)
		return nil, err
	}

	candidates := findCandidatesDirectories(patterns)
	watchDirs := findActualDirs(candidates, maxWatchDirs)

	if len(watchDirs) > maxWatchDirs {
		logger.Error(strings.Join(watchDirs[:maxWatchDirs], "\n") + "\n...")
		return nil, errors.New("too many watchDirs")
	}

	for _, t := range watchDirs {
		err = watcher.Add(t)
		if err != nil {
			if err.Error() == "bad file descriptor" {
				logger.Info("%s: %v", t, err)
			} else {
				logger.Error("%s: %v", t, err)
			}
		} else {
			logger.Info("gazing at: %s", t)
		}
	}

	notify := &Notify{
		Events:                  make(chan Event),
		watcher:                 watcher,
		isClosed:                false,
		times:                   make(map[string]int64),
		pendingPeriod:           100,
		regardRenameAsModPeriod: 1000,
		detectCreate:            true,
		candidates:              candidates,
	}

	go notify.wait()

	return notify, nil
}

func findActualDirs(patterns []string, maxWatchDirs int) []string {
	targets := uniq.New()

	for _, pattern := range patterns {
		dirs := findDirsByPattern(pattern)
		targets.AddAll(dirs)

		if targets.Len() > maxWatchDirs {
			break
		}
	}
	return targets.List()
}

// ["aaa/bbb/ccc"] -> [".", "aaa", "aaa/bbb", "aaa/bbb/ccc"]
// ["../aaa/bbb/ccc"] -> ["..", "../aaa", "../aaa/bbb", "../aaa/bbb/ccc"]
// ["/aaa/bbb/ccc"] -> ["/", "/aaa", "/aaa/bbb", "/aaa/bbb/ccc"]
func findCandidatesDirectories(patterns []string) []string {
	targets := uniq.New()

	for _, pattern := range patterns {
		paths := parsePathPattern(pattern)
		for i := len(paths) - 1; i >= 0; i-- {
			targets.Add(paths[i])
		}
	}
	return targets.List()
}

// "aaa/bbb/ccc/*/ddd/eee/*" -> ["aaa/bbb/ccc/*/ddd/eee/*", "aaa/bbb/ccc/*/ddd/eee", "aaa/bbb/ccc/*/ddd", "aaa/bbb/ccc/*", "aaa/bbb/ccc", "aaa/bbb", "aaa", "."]
func parsePathPattern(pathPattern string) []string {
	result := []string{}

	if len(pathPattern) == 0 {
		return result
	}
	if pathPattern == "/" || pathPattern == "\\" || pathPattern == "." || pathPattern == ".." {
		return []string{pathPattern}
	}

	result = append(result, pathPattern)

	isAbs := filepath.IsAbs(pathPattern) || pathPattern[0] == '/'
	isParent := strings.HasPrefix(pathPattern, "..")
	isWinAbs := pathPattern[0] != '/' && isAbs
	isExplicitCurrent := false
	isCurrent := !isAbs && !isParent
	if isCurrent {
		isExplicitCurrent = strings.HasPrefix(pathPattern, ".")
	}

	winFirstDelimiter := -1
	if isWinAbs {
		winFirstDelimiter = strings.Index(pathPattern, "\\")
	}

	for i := len(pathPattern) - 1; i >= 0; i-- {
		ch := pathPattern[i]

		if ch == '/' || ch == '\\' {

			if i > 0 {
				p := pathPattern[0:i]
				if winFirstDelimiter == i {
					p += "\\"
				}
				result = append(result, p)
			} else {
				result = append(result, "/")
			}
		}
	}

	if len(result) <= 1 {
		if !isAbs && !isExplicitCurrent {
			result = append(result, ".")
		}
	} else {
		if isCurrent && !isExplicitCurrent {
			result = append(result, ".")
		}
	}

	return result
}

func findDirsByPattern(pattern string) []string {
	patternDir := filepath.Dir(pattern)
	logger.Debug("pattern: %s", pattern)
	logger.Debug("patternDir: %s", patternDir)

	var targets []string

	realDir := findRealDirectory(patternDir)
	if len(realDir) > 0 {
		targets = append(targets, realDir)
	}

	_, dirs1 := fs.Find(pattern)
	targets = append(targets, dirs1...)

	_, dirs2 := fs.Find(patternDir)
	targets = append(targets, dirs2...)

	return targets
}

func findRealDirectory(path string) string {
	entries := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")

	currentPath := ""
	for i := 0; i < len(entries); i++ {
		if containsWildcard(entries[i]) {
			break
		}

		currentPath += entries[i] + string(filepath.Separator)
	}
	currentPath = fs.TrimSuffix(currentPath, string(filepath.Separator))

	if fs.IsDir(currentPath) {
		return currentPath
	} else {
		return ""
	}
}

func containsWildcard(path string) bool {
	return strings.ContainsAny(path, "*?[{")
}

func shouldWatch(dirPath string, candidates []string) bool {
	dirPathSlash := filepath.ToSlash(dirPath)
	if !fs.IsDir(dirPathSlash) {
		return false
	}

	for _, pattern := range candidates {
		patternSlash := filepath.ToSlash(pattern)
		ok, _ := doublestar.Match(patternSlash, dirPathSlash)
		if ok {
			return true
		}
	}

	return false
}

func (n *Notify) watchNewDirRecursive(dirPath string) {
	n.watchNewDir(dirPath)
	subDirs, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Error("ReadDir: %s", err)
		return
	}
	for _, subDir := range subDirs {
		if subDir.IsDir() {
			subDirPath := filepath.Join(dirPath, subDir.Name())
			n.watchNewDirRecursive(subDirPath)
		}
	}
}

func (n *Notify) wait() {
	for {
		select {
		case event, ok := <-n.watcher.Events:
			normalizedName := filepath.Clean(event.Name)

			logger.Debug("fs.IsDir: %s", fs.IsDir(normalizedName))
			if event.Has(fsnotify.Create) && shouldWatch(normalizedName, n.candidates) {
				logger.Info("gazing at: %s", normalizedName)
				n.watchNewDirRecursive(normalizedName)
			}

			if !ok {
				continue
			}
			if !n.shouldExecute(normalizedName, event) {
				continue
			}
			logger.Debug("notified: %s: %s", normalizedName, event.Op)
			now := time.UnixNano()
			n.times[normalizedName] = now
			e := Event{
				Name: normalizedName,
				Time: now,
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

func (n *Notify) watchNewDir(normalizedName string) {
	err := n.watcher.Remove(normalizedName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "fsnotify: can't remove non-existent") {
			logger.Debug("watcher.Remove: %s", err)
		} else {
			logger.Error("watcher.Remove: %s", err)
		}
	}
	err = n.watcher.Add(normalizedName)
	if err != nil {
		logger.Error("watcher.Add: %s", err)
	}
}

func (n *Notify) shouldExecute(filePath string, ev fsnotify.Event) bool {
	const W = fsnotify.Write
	const R = fsnotify.Rename
	const C = fsnotify.Create

	if !ev.Has(W) && !ev.Has(R) && !(n.detectCreate && ev.Has(C)) {
		logger.Debug("skipped: %s: %s (Op is not applicable)", filePath, ev.Op)
		return false
	}

	lastExecutionTime := n.times[filePath]

	if !fs.IsFile(filePath) {
		logger.Debug("skipped: %s: %s (not a file)", filePath, ev.Op)
		return false
	}

	if strings.Contains(filePath, "'") || strings.Contains(filePath, "\"") {
		logger.Debug("skipped: %s: %s (unsupported character)", filePath, ev.Op)
		return false
	}

	modifiedTime := time.GetFileModifiedTime(filePath)

	if ev.Has(W) || ev.Has(C) {
		elapsed := modifiedTime - lastExecutionTime
		logger.Debug("lastExecutionTime(%s): %d, %d", ev.Op, lastExecutionTime, elapsed)
		if elapsed < n.pendingPeriod*1000000 {
			logger.Debug("skipped: %s: %s (too frequent)", filePath, ev.Op)
			return false
		}
	}
	if ev.Has(R) {
		elapsed := time.UnixNano() - modifiedTime
		logger.Debug("lastExecutionTime(%s): %d, %d", ev.Op, lastExecutionTime, elapsed)
		if elapsed > n.regardRenameAsModPeriod*1000000 {
			logger.Debug("skipped: %s: %s (unnatural rename)", filePath, ev.Op)
			return false
		}
	}

	return true
}

// PendingPeriod sets new pendingPeriod(ms).
func (n *Notify) PendingPeriod(p int64) {
	n.pendingPeriod = p
}

// Requeue requeue an event.
func (n *Notify) Requeue(event Event) {
	n.Events <- event
}
