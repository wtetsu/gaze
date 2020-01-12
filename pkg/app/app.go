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
func Start(watchFiles []string, userCommand string, file string, timeout int) error {
	watcher := gazer.New(watchFiles)
	defer watcher.Close()

	commandConfigs, err := createCommandConfig(userCommand, file)
	if err != nil {
		return err
	}
	err = watcher.Run(commandConfigs, timeout)
	return err
}

func createCommandConfig(userCommand string, file string) (*config.Config, error) {
	if userCommand != "" {
		logger.Debug("userCommand: %s", userCommand)
		commandConfigs := config.New(userCommand)
		return commandConfigs, nil
	}

	if file != "" {
		return config.LoadConfig(file)
	}

	return config.InitConfig()
}
