/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/mattn/go-shellwords"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

func executeCommandOrTimeout(cmd *exec.Cmd, timeout <-chan struct{}) error {
	exec := executeCommandAsync(cmd)

	var err error
	finished := false
	for {
		if finished {
			break
		}
		select {
		case <-timeout:
			if cmd.Process == nil {
				timeout = time.After(5)
				continue
			}
			kill(cmd, "Timeout")
			finished = true
			err = errors.New("")
		case err = <-exec:
			finished = true
		}
	}
	if err != nil {
		return err
	}

	if cmd.ProcessState != nil {
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			return fmt.Errorf("exitCode:%d", exitCode)
		}
	}

	return nil
}

func executeCommandAsync(cmd *exec.Cmd) <-chan error {
	ch := make(chan error)

	go func() {
		if cmd == nil {
			ch <- errors.New("failed: cmd is nil")
			return
		}
		err := executeCommand(cmd)
		if err != nil {
			ch <- err
			return
		}
		ch <- nil
	}()
	return ch
}

func executeCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	if cmd.Process != nil {
		logger.Info("Pid: %d", cmd.Process.Pid)
	} else {
		logger.Info("Pid: ????")
	}
	err := cmd.Wait()
	return err
}

func kill(cmd *exec.Cmd, reason string) bool {
	if cmd == nil || cmd.Process == nil {
		return false
	}
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		return false
	}

	var signal os.Signal
	if runtime.GOOS == "windows" {
		signal = os.Kill
	} else {
		signal = syscall.SIGTERM
	}
	err := cmd.Process.Signal(signal)
	if err != nil {
		logger.Notice("kill failed: %v", err)
		return false
	}
	logger.Notice("%s: %d has been killed", reason, cmd.Process.Pid)
	return true
}

func createCommand(commandString string) *exec.Cmd {
	parser := shellwords.NewParser()
	// parser.ParseBacktick = true
	// parser.ParseEnv = true
	args, err := parser.Parse(commandString)
	if err != nil {
		return nil
	}
	if len(args) == 1 {
		return exec.Command(args[0])
	}
	return exec.Command(args[0], args[1:]...)
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
