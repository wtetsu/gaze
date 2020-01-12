package gazer

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-shellwords"
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
func (g *Gazer) Run(configs *config.Config, timeout int) error {
	err := waitAndRunForever(g.watcher, g.patterns, configs, timeout)
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

func waitAndRunForever(watcher *fsnotify.Watcher, watchFiles []string, commandConfigs *config.Config, timeout int) error {
	var lastExecutionTime int64

	sigInt := sigIntChannel()

	isDisposed := false
	for {
		if isDisposed {
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
				logger.NoticeWithBlank("[%s]", commandString)

				err := executeCommandOrTimeout(commandString, timeout)
				if err != nil {
					logger.NoticeObject(err)
				}
			}
			lastExecutionTime = time.Now()
		case <-sigInt:
			isDisposed = true
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

func executeCommandOrTimeout(commandString string, timeoutMill int) error {
	cmd, exec := executeCommand(commandString)
	if timeoutMill <= 0 {
		return <-exec
	}

	timeout := time.After(timeoutMill)
	var err error

	finished := false
	for {
		if finished {
			break
		}
		select {
		case <-timeout:
			if cmd.Process == nil {
				timeout = time.After(timeoutMill)
				break
			}
			err = cmd.Process.Kill()
			if err != nil {
				logger.NoticeObject(err)
			}
			logger.Notice("Timeout: %d has been killed", cmd.Process.Pid)
			finished = true

		case err = <-exec:
			finished = true
		}
	}
	return err
}

func executeCommand(commandString string) (*exec.Cmd, <-chan error) {
	ch := make(chan error)
	cmd := createCommand(commandString)

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

func createCommand(commandString string) *exec.Cmd {
	args, err := shellwords.Parse(commandString)
	if err != nil {
		return nil
	}
	if len(args) == 1 {
		return exec.Command(args[0])
	}
	return exec.Command(args[0], args[1:]...)
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
