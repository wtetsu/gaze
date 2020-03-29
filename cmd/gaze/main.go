/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package main

import (
	"fmt"
	"os"

	"github.com/wtetsu/gaze/pkg/app"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
)

const version = "v0.1.5"

func main() {
	args := app.ParseArgs(os.Args, func() {
		fmt.Println(usage2())
	})

	if args.Help() {
		fmt.Println(usage2())
		return
	}

	if args.Version() {
		fmt.Println("gaze " + version)
		return
	}

	if args.Yaml() {
		fmt.Println(config.Default())
		return
	}

	if len(args.Targets()) == 0 {
		fmt.Println(usage1())
		return
	}

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

	err := app.Start(args.Targets(), args.UserCommand(), args.File(), args.Timeout(), args.Restart())
	if err != nil {
		logger.ErrorObject(err)
		os.Exit(1)
	}

}

func usage1() string {
	return `Usage: gaze [options...] file(s)

Options(excerpt):
  -c  A command string.
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -h  Display help.

Examples:
  gaze .
  gaze main.go
  gaze a.rb b.rb
  gaze -c make '**/*.c'
  gaze -c "eslint {{file}}" 'src/**/*.js'
  gaze -r server.py
  gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze`
}

func usage2() string {
	return `Usage: gaze [options...] file(s)

Options:
  -c  A command string.
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -f  Specify a YAML configuration file.
  -v  Verbose mode.
  -q  Quiet mode.
  -y  Output the default configuration.
  -h  Display help.
  --color    Color mode (0:plain, 1:colorful).
  --version  Output version information.

Examples:
  gaze .
  gaze main.go
  gaze a.rb b.rb
  gaze -c make '**/*.c'
  gaze -c "eslint {{file}}" 'src/**/*.js'
  gaze -r server.py
  gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze`
}
