package gazer

import (
	"errors"
	"fmt"
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
	isClosed bool
}

// New returns a new Gazer.
func New(patterns []string) *Gazer {
	watcher, _ := createWatcher(patterns)
	return &Gazer{
		patterns: patterns,
		watcher:  watcher,
		isClosed: false,
	}
}

// Close disposes internal resources.
func (g *Gazer) Close() {
	if g.isClosed {
		return
	}
	g.watcher.Close()
	g.isClosed = true
}

// Run starts to gaze.
func (g *Gazer) Run(configs *config.Config) error {
	err := waitAndRunForever(g.watcher, g.patterns, configs)
	return err
}

func createWatcher(patterns []string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.ErrorObject(err)
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
			logger.Info("gazing at: %s", d)
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
				logger.Notice("[%s]", commandString)

				err := executeShellCommandOrTimeout(cmd.Shell(), scriptPath, 1000)
				if err != nil {
					logger.NoticeObject(err)
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

func executeShellCommandOrTimeout(shell string, scriptPath string, timeoutMill int) error {
	timeout := time.After(timeoutMill)
	cmd, exec := executeShellCommand(shell, scriptPath)

	var err error
	select {
	case <-timeout:
		fmt.Println(cmd.Process.Pid)
		if cmd.Process != nil {
			cmd.Process.Kill()
			fmt.Println("kill!!!")
		}
	case err = <-exec:
		// if cmd != nil {
		// 	if cmd.ProcessState != nil {
		// 		logger.Info("exit: %d", cmd.ProcessState.ExitCode())
		// 	}
		// }
	}

	return err
}

func executeShellCommand(shell string, scriptPath string) (*exec.Cmd, <-chan error) {
	ch := make(chan error)
	cmd := createScriptCommand(shell, scriptPath)

	go func() {
		if cmd == nil {
			ch <- errors.New("failed")
			return
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		logger.Info("Pid: %d", cmd.Process.Pid)

		err := cmd.Wait()
		if err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()
	return cmd, ch
}

func createScriptCommand(shell string, scriptPath string) *exec.Cmd {
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
