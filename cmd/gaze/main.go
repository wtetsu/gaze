/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wtetsu/gaze/pkg/app"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
)

//go:embed version
var version string

const (
	errTimeout      = "timeout must be more than 0"
	errColor        = "color must be 0 or 1"
	errMaxWatchDirs = "maxWatchDirs must be more than 0"
)

func main() {
	args := app.ParseArgs(os.Args, func() {
		fmt.Println(usage2())
	})

	if !args.Debug() {
		// panic handler
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorObject(err)
				os.Exit(1)
			}
		}()
	}

	done, exitCode := earlyExit(args)
	if done {
		os.Exit(exitCode)
		return
	}

	initLogger(args)

	err := validate(args)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	appOptions := app.NewAppOptions(args.Timeout(), args.Restart(), args.MaxWatchDirs())

	err = app.Start(args.Targets(), args.UserCommand(), args.File(), appOptions)
	if err != nil {
		logger.ErrorObject(err)
		os.Exit(1)
	}
}

func earlyExit(args *app.Args) (bool, int) {
	if args.Help() {
		fmt.Println(usage2())
		return true, 0
	}

	if args.Version() {
		fmt.Println("gaze " + version)
		return true, 0
	}

	if args.Yaml() {
		fmt.Println(config.Default())
		return true, 0
	}

	if len(args.Targets()) == 0 {
		fmt.Println(usage1())
		return true, 1
	}
	return false, 0
}

func initLogger(args *app.Args) {
	if args.Color() == 0 {
		logger.Plain()
	} else {
		logger.Colorful()
	}
	if args.Quiet() {
		logger.Level(logger.QUIET)
	}
	if args.Verbose() {
		logger.Level(logger.VERBOSE)
	}
	if args.Debug() {
		logger.Level(logger.DEBUG)
	}
}

func validate(args *app.Args) error {
	var errorList []string
	if args.Timeout() <= 0 {
		errorList = append(errorList, errTimeout)
	}
	if args.Color() != 0 && args.Color() != 1 {
		errorList = append(errorList, errColor)
	}
	if args.MaxWatchDirs() <= 0 {
		errorList = append(errorList, errMaxWatchDirs)
	}
	if len(errorList) >= 1 {
		return errors.New(strings.Join(errorList, "\n"))
	}
	return nil
}

func usage1() string {
	return `Usage: gaze [options] file(s)

Options(excerpt):
  -c <command>    Command(s) to run when files change.
  -r              Restart mode: send SIGTERM to the running process before starting the next command.
  -t <time_ms>    Timeout (ms): send SIGTERM to the running process after the specified time.
  -h              Show help.

Examples:
  gaze .
  gaze main.go
  gaze a.rb b.rb
  gaze -c make "**/*.c"
  gaze -c "eslint {{file}}" "src/**/*.js"
  gaze -r server.py
  gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze`
}

func usage2() string {
	return `Usage: gaze [options] file(s)

Options:
  -c <command>    Command(s) to run when files change.
  -r              Restart mode: send SIGTERM to the running process before starting the next command.
  -t <time_ms>    Timeout (ms): send SIGTERM to the running process after the specified time.
  -f <file>       Path to a YAML configuration file.
  -v              Verbose mode: show additional information.
  -q              Quiet mode: suppress normal output.
  -y              Show the default YAML configuration.
  -h              Show help.
  --color <mode>  Set color mode (0: plain, 1: colorful).
  --version       Show version information.

Examples:
  gaze .
  gaze main.go
  gaze a.rb b.rb
  gaze -c make "**/*.c"
  gaze -c "eslint {{file}}" "src/**/*.js"
  gaze -r server.py
  gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze`
}
