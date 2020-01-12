/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package app

import (
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/gazer"
	"github.com/wtetsu/gaze/pkg/logger"
)

// Start starts a gaze process
func Start(watchFiles []string, userCommand string, timeout int) error {
	watcher := gazer.New(watchFiles)
	defer watcher.Close()

	commandConfigs, err := createCommandConfig(userCommand)
	if err != nil {
		return err
	}
	err = watcher.Run(commandConfigs, timeout)
	return err
}

func createCommandConfig(userCommand string) (*config.Config, error) {
	if userCommand != "" {
		logger.Debug("userCommand: %s", userCommand)
		commandConfigs := config.New(userCommand)
		return commandConfigs, nil
	}

	commandConfigs, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return commandConfigs, nil
}
