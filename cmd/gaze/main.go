/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wtetsu/gaze/pkg/app"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
)

const version = "v1.1.7"

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
		errorList = append(errorList, "timeout must be more than 0")
	}
	if args.Color() != 0 && args.Color() != 1 {
		errorList = append(errorList, "color must be 0 or 1")
	}
	if args.MaxWatchDirs() <= 0 {
		errorList = append(errorList, "maxWatchDirs must be more than 0")
	}
	if len(errorList) >= 1 {
		return errors.New(strings.Join(errorList, "\n"))
	}
	return nil
}

func usage1() string {
	return `Usage: gaze [options...] file(s)

Options(excerpt):
  -c  Command(s).
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -h  Display help.

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
	return `Usage: gaze [options...] file(s)

Options:
  -c  Command(s).
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -f  Specify a YAML configuration file.
  -v  Verbose mode.
  -q  Quiet mode.
  -y  Display the default YAML configuration.
  -h  Display help.
  --color    Color mode (0:plain, 1:colorful).
  --version  Display version information.

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
