package app

import (
	"github.com/wtetsu/gaze/pkg/proc"
)

// Start starts a gaze process
func Start(files []string, userCommand string) error {
	err := proc.StartGazing(files, userCommand)
	return err
}
