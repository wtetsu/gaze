package proc

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/file"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

var commandFileMap = make(map[string]string)

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

func executeShellCommand(commandString string) error {
	defaultShell := getDefaultShell()

	shellScriptPath := prepareScript(defaultShell, commandString)
	cmd := executeScript(defaultShell, shellScriptPath)

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

func prepareScript(defaultShell string, commandString string) string {
	existingFilePath, found := commandFileMap[commandString]

	if found && file.Exist(existingFilePath) {
		return existingFilePath
	}

	newFilePath, err := ioutil.TempFile("", "*.gaze.cmd")
	if err != nil {
		return ""
	}

	if defaultShell == "cmd" {
		newFilePath.WriteString("@" + commandString)
	} else {
		newFilePath.WriteString(commandString)
	}
	err = newFilePath.Close()
	if err != nil {
		return ""
	}

	commandFileMap[commandString] = newFilePath.Name()

	return newFilePath.Name()
}

func executeScript(shell string, scriptPath string) *exec.Cmd {
	if shell == "cmd" {
		return exec.Command("cmd", "/c", scriptPath)
	}
	return exec.Command(shell, scriptPath)
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
