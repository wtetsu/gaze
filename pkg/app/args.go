/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package app

// Args has application arguments
type Args struct {
	help        bool
	restart     bool
	userCommand string
	timeout     int
	yaml        bool
	quiet       bool
	verbose     bool
	file        string
	color       int
	debug       bool
	version     bool
	targets     []string
}

// Help returns a.help
func (a *Args) Help() bool {
	return a.help
}

// Restart returns a.restart
func (a *Args) Restart() bool {
	return a.restart
}

// UserCommand returns a.userCommand
func (a *Args) UserCommand() string {
	return a.userCommand
}

// Timeout returns a.timeout
func (a *Args) Timeout() int {
	return a.timeout
}

// Yaml returns a.yaml
func (a *Args) Yaml() bool {
	return a.yaml
}

// Quiet returns a.quiet
func (a *Args) Quiet() bool {
	return a.quiet
}

// Verbose returns a.verbose
func (a *Args) Verbose() bool {
	return a.verbose
}

// File returns a.file
func (a *Args) File() string {
	return a.file
}

// Color returns a.color
func (a *Args) Color() int {
	return a.color
}

// Debug returns a.debug
func (a *Args) Debug() bool {
	return a.debug
}

// Version returns a.version
func (a *Args) Version() bool {
	return a.version
}

// Targets returns a.targets
func (a *Args) Targets() []string {
	return a.targets
}
