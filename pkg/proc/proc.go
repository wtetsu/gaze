package proc

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/fsnotify/fsnotify"
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

func waitAndRunForever(watcher *fsnotify.Watcher, files []string, userCommand string) {
	var lastExecutionTime int64

	commandConfigs := createCommandConfigs()

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
			if userCommand != "" {
				executeShellCommand(userCommand)
			} else {
				command := createCommand(event.Name, commandConfigs)
				if command != nil {
					executeCommand(command)
				}
			}

			lastExecutionTime = time.Now()
		}
	}
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

func executeCommand(command []string) {
	cmd := exec.Command(command[0], command[1:]...)
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
}

func executeShellCommand(commandString string) {
	tmpFile, err := ioutil.TempFile("", "*.sh")
	if err != nil {
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(commandString)
	if err != nil {
		return
	}

	var cmd *exec.Cmd
	shell := os.Getenv("SHELL")
	if shell != "" {
		cmd = exec.Command(shell, tmpFile.Name())
	} else {
		cmd = exec.Command("sh", tmpFile.Name())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	// fmt.Println(cmd.Process.Pid)
	// time.Sleep(2000)
	// cmd.Process.Kill()

	err = cmd.Wait()
	if err != nil {
		logger.Debug(err)
	}
}

func createCommand(filePath string, configs []commandConfig) []string {
	var command []string = nil
	for _, c := range configs {
		if c.SearchRegexp == nil {
			continue
		}
		if c.SearchRegexp.MatchString(filePath) {
			command = append(c.Command, filePath)
			break
		}
	}
	return command
}

func createCommandConfigs() []commandConfig {
	var resultConfigs []commandConfig
	configs := getConfigs()

	n := len(configs)
	for i := 0; i < n; i++ {
		c := configs[i]
		re, err := regexp.Compile(c.Search)
		if err != nil {
			logger.Fatal(err)
			continue
		}
		c.SearchRegexp = re
		resultConfigs = append(resultConfigs, commandConfig{c.Search, c.Command, re})
	}

	return resultConfigs
}

// TODO
func getConfigs() []commandConfig {
	return []commandConfig{
		commandConfig{`\.d`, []string{"dmd", "-run"}, nil},
		commandConfig{`\.js`, []string{"node"}, nil},
		commandConfig{`\.go`, []string{"go run"}, nil},
		commandConfig{`\.php`, []string{"php"}, nil},
		commandConfig{`\.pl`, []string{"perl"}, nil},
		commandConfig{`\.py`, []string{"python"}, nil},
		commandConfig{`\.rb`, []string{"ruby"}, nil},
		commandConfig{`\.sh`, []string{"sh"}, nil},
	}

}

type commandConfig struct {
	Search       string
	Command      []string
	SearchRegexp *regexp.Regexp
}
