/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package app

type AppOptions struct {
	timeout      int64
	restart      bool
	maxWatchDirs int
}

func NewAppOptions(timeout int64, restart bool, maxWatchDirs int) AppOptions {
	return AppOptions{
		timeout:      timeout,
		restart:      restart,
		maxWatchDirs: maxWatchDirs,
	}
}

func (a AppOptions) Timeout() int64 {
	return a.timeout
}

func (a AppOptions) Restart() bool {
	return a.restart
}

func (a AppOptions) MaxWatchDirs() int {
	return a.maxWatchDirs
}
