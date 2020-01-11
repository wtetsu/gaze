package gazer

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/wtetsu/gaze/pkg/command"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

// Gazer gazes filesystem.
type Gazer struct {
	patterns []string
	watcher  *fsnotify.Watcher
}

// New returns a new Gazer.
func New(patterns []string) *Gazer {
	watcher, _ := createWatcher(patterns)
	return &Gazer{
		patterns: patterns,
		watcher:  watcher,
	}
}

// Run starts to gaze.
func (g *Gazer) Run(configs *config.Config) error {
	err := waitAndRunForever(g.watcher, g.patterns, configs)
	return err
}

func createWatcher(patterns []string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	added := map[string]struct{}{}
	for _, pattern := range patterns {
		_, dirs := fs.Find(pattern)
		for _, d := range dirs {
			_, ok := added[d]
			if ok {
				continue
			}
			logger.Notice("gazing at: %s", d)
			err = watcher.Add(d)
			if err != nil {
				logger.DebugObject(err)
			}
			added[d] = struct{}{}
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

				logger.Notice(commandString)

				err := executeShellCommand(cmd.Shell(), scriptPath)
				if err != nil {
					logger.Error(err)
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
		if fs.GlobMatch(f, s) {
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
		logger.DebugObject(err)
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
