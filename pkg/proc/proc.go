package proc

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

func StartWatcher(command []string, files []string) {
	commandString := strings.Join(command[:], " ")
	logger.Debug(commandString)

	watcher, err := createWatcher(files)
	if err != nil {
		return
	}

	defer watcher.Close()

	done := make(chan bool)

	go waitAndRunForever(command, watcher)

	<-done
}

func createWatcher(files []string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}
	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			logger.Fatal(err)
		}
	}
	return watcher, nil
}

func waitAndRunForever(command []string, watcher *fsnotify.Watcher) {
	var lastExecutionTime int64
	for {
		select {
		case event, ok := <-watcher.Events:
			if ok && event.Op&fsnotify.Write == fsnotify.Write {
				modifiedTime := time.GetFileModifiedTime(event.Name)
				if modifiedTime > lastExecutionTime {
					executeCommand(command)
					lastExecutionTime = time.Now()
				}
			}
		}
	}
}

func executeCommand(command []string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	// fmt.Println(cmd.Process.Pid)
	// time.Sleep(2 * time.Second)
	// cmd.Process.Kill()

	err := cmd.Wait()
	if err != nil {
		os.Exit(1)
	}
}
