package proc

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/bmatcuk/doublestar"
	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/command"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
	"github.com/wtetsu/gaze/pkg/fs"
)

// StartGazing starts file
func StartGazing(watchFiles []string, userCommand string) error {
	watcher, err := createWatcher(watchFiles)
	if err != nil {
		return err
	}
	defer watcher.Close()

	var commandConfigs *config.Config
	if userCommand != "" {
		logger.Debugf("userCommand: %s", userCommand)
		commandConfigs = config.New(userCommand)
	} else {
		commandConfigs, err = config.LoadConfig()
		if err != nil {
			return err
		}
	}

	logger.Debug(commandConfigs)

	err = waitAndRunForever(watcher, watchFiles, commandConfigs)

	return err
}

func createWatcher(watchFiles []string) (*fsnotify.Watcher, error) {
	logger.Debug("createWatcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	dirs := fs.ListDir(".", "**")
	for _, d := range dirs {
		logger.Debugf("watching: %s", d)
		err = watcher.Add(d)
		if err != nil {
			logger.Debug(err)
		}
	}

	return watcher, nil
}

func waitAndRunForever(watcher *fsnotify.Watcher, watchFiles []string, commandConfigs *config.Config) error {
	cmd := command.New(getDefaultShell())
	defer cmd.Dispose()

	var lastExecutionTime int64

	sigInt := sigIntChannel()

	for {
		if cmd.Disposed() {
			break
		}
		select {
		case event, ok := <-watcher.Events:
			flag := fsnotify.Write | fsnotify.Rename
			if ok && event.Op|flag == 0 {
				continue
			}
			if !matchAny(watchFiles, event.Name) {
				continue
			}
			modifiedTime := time.GetFileModifiedTime(event.Name)
			if modifiedTime <= lastExecutionTime {
				continue
			}

			commandString := getAppropriateCommand(event.Name, commandConfigs)
			if commandString != "" {
				scriptPath := cmd.PrepareScript(commandString)
				fmt.Println(scriptPath)

				err := executeShellCommand(cmd.Shell(), scriptPath)
				if err != nil {
					logger.Fatal(err)
				}
			}
			lastExecutionTime = time.Now()
		case <-sigInt:
			cmd.Dispose()
			return nil
		}
	}
	return nil
}

func sigIntChannel() chan struct{} {
	ch := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(ch)
	}()
	return ch
}

func matchAny(watchFiles []string, s string) bool {
	result := false
	for _, f := range watchFiles {
		ok, _ := doublestar.Match(f, s)
		if ok {
			result = true
			break
		}
	}
	return result
}

func getDefaultShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return shell
	}
	if runtime.GOOS == "windows" {
		return "cmd"
	}
	return "sh"
}

func executeShellCommand(shell string, scriptPath string) error {
	cmd := executeScript(shell, scriptPath)

	if cmd == nil {
		return errors.New("failed")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	// fmt.Println(cmd.Process.Pid)
	// time.Sleep(2000)
	// cmd.Process.Kill()

	err := cmd.Wait()
	if err != nil {
		logger.Debug(err)
	}
	return nil
}

func executeScript(shell string, scriptPath string) *exec.Cmd {
	if shell == "cmd" {
		return exec.Command("cmd", "/c", scriptPath)
	}
	return exec.Command(shell, scriptPath)
}

func getAppropriateCommand(filePath string, commandConfigs *config.Config) string {
	var result string
	for _, c := range commandConfigs.Commands {
		if c.Run == "" || c.Ext == "" && c.Re == "" {
			continue
		}
		if c.Match(filePath) {
			command := render(c.Run, filePath)
			result = command
			break
		}
	}
	return result
}
