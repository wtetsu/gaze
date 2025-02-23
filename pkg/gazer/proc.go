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
	"time"

	"github.com/mattn/go-shellwords"
	"github.com/wtetsu/gaze/pkg/gutil"
	"github.com/wtetsu/gaze/pkg/logger"
)

type CmdResult struct {
	StartTime time.Time
	EndTime   time.Time
	Err       error
}

func executeCommandOrTimeout(cmd *exec.Cmd, timeoutMills int64) CmdResult {
	exec := executeCommandAsync(cmd)

	var cmdResult CmdResult
	var launchedTime = time.Now()
	finished := false
	timeout := gutil.After(timeoutMills)
	for {
		if finished {
			break
		}
		select {
		case <-timeout:
			if cmd.Process == nil {
				timeout = gutil.After(5)
				continue
			}
			kill(cmd, "Timeout")
			finished = true
			cmdResult = CmdResult{StartTime: launchedTime, EndTime: time.Now(), Err: errors.New("")}
		case cmdResult = <-exec:
			finished = true
		}
	}
	if cmdResult.Err != nil {
		return cmdResult
	}

	if cmd.ProcessState != nil {
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode != 0 {
			cmdResult.Err = fmt.Errorf("exitCode:%d", exitCode)
			return cmdResult
		}
	}

	return cmdResult
}

func executeCommandAsync(cmd *exec.Cmd) <-chan CmdResult {
	ch := make(chan CmdResult)

	go func() {
		if cmd == nil {
			ch <- CmdResult{Err: errors.New("failed: cmd is nil")}
			return
		}
		cmdResult := executeCommand(cmd)
		ch <- cmdResult
	}()
	return ch
}

func executeCommand(cmd *exec.Cmd) CmdResult {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		return CmdResult{StartTime: start, EndTime: time.Now(), Err: err}
	}

	if cmd.Process != nil {
		logger.Info("Pid: %d", cmd.Process.Pid)
	} else {
		logger.Info("Pid: ????")
	}
	err = cmd.Wait()

	return CmdResult{StartTime: start, EndTime: time.Now(), Err: err}
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
