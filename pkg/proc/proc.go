package proc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

// StartGazing starts file
func StartGazing(files []string, userCommand string) error {
	watcher, err := createWatcher(files)
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)

	go waitAndRunForever(watcher, files, userCommand)

	<-done

	return nil
}

func createWatcher(files []string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}
	err = watcher.Add(".")
	if err != nil {
		logger.Debug(err)
	}
	return watcher, nil
}

func waitAndRunForever(watcher *fsnotify.Watcher, files []string, userCommand string) error {
	var lastExecutionTime int64

	commandConfigs, err := config.LoadConfig()
	if err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if ok && event.Op&fsnotify.Write != fsnotify.Write {
				continue
			}
			if !match(files, event.Name) {
				continue
			}
			modifiedTime := time.GetFileModifiedTime(event.Name)
			if modifiedTime <= lastExecutionTime {
				continue
			}
			var command string
			if userCommand != "" {
				command = userCommand
			} else {
				command = getAppropriateCommand(event.Name, commandConfigs)
			}
			if command != "" {
				err := executeShellCommand(command)
				if err != nil {
					logger.Fatal(err)
				}
			}
			lastExecutionTime = time.Now()
		}
	}
	// return nil
}

func match(files []string, s string) bool {
	result := false
	for _, f := range files {
		if f == s {
			result = true
			break
		}
	}
	return result
}

func executeShellCommand(commandString string) error {
	var cmd *exec.Cmd

	shell := os.Getenv("SHELL")
	if shell != "" {
		cmd = executeSh(shell, commandString)
	} else {
		if runtime.GOOS == "windows" {
			cmd = executeBat(commandString)
		} else {
			cmd = executeSh("sh", commandString)
		}
	}

	if cmd == nil {
		return errors.New("failed:" + commandString)
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

func executeSh(shell string, commandString string) *exec.Cmd {
	tmpFile, err := ioutil.TempFile("", "*.sh")
	if err != nil {
		return nil
	}

	fmt.Println(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(commandString)
	if err != nil {
		return nil
	}

	return exec.Command(shell, tmpFile.Name())
}

func executeBat(commandString string) *exec.Cmd {
	tmpFile, err := ioutil.TempFile("", "*.bat")
	if err != nil {
		return nil
	}

	fmt.Println(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("@" + commandString)
	if err != nil {
		return nil
	}

	return exec.Command("cmd", "/c", tmpFile.Name())
}

func getAppropriateCommand(filePath string, commandConfigs []config.Config) string {
	ext := filepath.Ext(filePath)
	// base := filepath.Base(filePath)
	// abs, _ := filepath.Abs(filePath)
	// dir := filepath.Dir(filePath)

	var result string
	for _, c := range commandConfigs {
		// if c.SearchRegexp.MatchString(filePath) {
		// 	command = append(c.Command, filePath)
		// 	break
		// }
		if c.Run != "" && c.Ext == ext {
			command := render(c.Run, filePath)
			result = command
			break
		}
	}
	return result
}
