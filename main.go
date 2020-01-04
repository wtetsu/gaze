package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/logger"
)

func main() {
	command := []string{"sleep", "3"}
	files := []string{"aaa.txt", "aaa2.txt"}
	startWatcher(command, files)
}

func startWatcher(command []string, files []string) {
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

func waitAndRunForever(command []string, watcher *fsnotify.Watcher) {
	for {
		waitAndRun(command, watcher)
	}
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

func waitAndRun(command []string, watcher *fsnotify.Watcher) {
	select {
	case event, ok := <-watcher.Events:
		if !ok {
			return
		}
		logger.Println("event:", event)
		if event.Op&fsnotify.Write == fsnotify.Write {
			// logger.Println("modified file:", event.Name)
			executeCommand(command)
		}
		// case err, ok := <-watcher.Errors:
		// 	if !ok {
		// 		return
		// 	}
		// 	logger.Println("error:", err)
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
