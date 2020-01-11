package proc

import (
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/gazer"
	"github.com/wtetsu/gaze/pkg/logger"
)

// StartGazing starts file
func StartGazing(watchFiles []string, userCommand string) error {
	// watcher, err := createWatcher(watchFiles)
	// if err != nil {
	// 	return err
	// }
	// defer watcher.Close()

	watcher := gazer.New(watchFiles)

	// logger.Debug(commandConfigs)
	// err = waitAndRunForever(watcher, watchFiles, commandConfigs)

	commandConfigs, err := createCommandConfig(userCommand)
	if err != nil {
		return err
	}
	err = watcher.Run(commandConfigs)
	return err
}

func createCommandConfig(userCommand string) (*config.Config, error) {
	if userCommand != "" {
		logger.Debugf("userCommand: %s", userCommand)
		commandConfigs := config.New(userCommand)
		return commandConfigs, nil
	}

	commandConfigs, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	return commandConfigs, nil

}
