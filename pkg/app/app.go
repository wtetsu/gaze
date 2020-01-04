package app

import (
	"path/filepath"
	"strings"

	"github.com/wtetsu/gaze/pkg/proc"
)

// Start starts a gaze process
func Start(args []string) {
	command, files := parseCommand(args)
	proc.StartWatcher(command, files)
}

func parseCommand(args []string) ([]string, []string) {
	config := map[string]string{
		".rb": "ruby",
		".py": "python",
		".js": "node",
		".d":  "dmd -run",
	}
	var command []string
	var files []string

	if len(args) == 1 {
		ext := filepath.Ext(args[0])
		exe := config[ext]

		command = append(command, exe, args[0])
		files = append(files, args[0])
	} else {
		command = strings.Split(args[0], " ")
		files = args[1:]
	}

	return command, files
}
