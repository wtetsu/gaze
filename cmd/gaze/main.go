/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wtetsu/gaze/pkg/app"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/uniq"
)

const version = "v0.0.8"

func main() {
	args := parseArgs()

	if args.help {
		fmt.Println(usage2())
		return
	}

	if args.version {
		fmt.Println("gaze " + version)
		return
	}

	if args.yaml {
		fmt.Println(config.Default())
		return
	}

	if len(args.targets) == 0 {
		fmt.Println(usage1())
		return
	}

	if args.color == 0 {
		logger.Plain()
	} else {
		logger.Colorful()
	}
	if args.quiet {
		logger.Level(logger.QUIET)
	}
	if args.verbose {
		logger.Level(logger.VERBOSE)
	}
	if args.debug {
		logger.Level(logger.DEBUG)
	}

	err := app.Start(args.targets, args.userCommand, args.file, args.timeout, args.restart)
	if err != nil {
		logger.ErrorObject(err)
		os.Exit(1)
	}

}

func parseArgs() *Args {
	flag.Usage = func() {
		usage1()
	}

	help := flag.Bool("h", false, "")
	restart := flag.Bool("r", false, "")
	userCommand := flag.String("c", "", "")
	timeout := flag.Int("t", 0, "")
	yaml := flag.Bool("y", false, "")
	quiet := flag.Bool("q", false, "")
	verbose := flag.Bool("v", false, "")
	file := flag.String("f", "", "")
	color := flag.Int("color", 1, "")
	debug := flag.Bool("debug", false, "")
	version := flag.Bool("version", false, "")

	files := []string{}
	optionStartIndex := len(os.Args)
	for i, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			optionStartIndex = i + 1
			break
		}
		files = append(files, a)
	}
	err := flag.CommandLine.Parse(os.Args[optionStartIndex:])
	if err != nil {
		return nil
	}

	u := uniq.New()
	u.AddAll(files)
	u.AddAll(flag.Args())

	args := Args{
		help:        *help,
		restart:     *restart,
		userCommand: *userCommand,
		timeout:     *timeout,
		yaml:        *yaml,
		quiet:       *quiet,
		verbose:     *verbose,
		debug:       *debug,
		file:        *file,
		color:       *color,
		version:     *version,
		targets:     u.List(),
	}

	return &args
}

// Args has application arguments
type Args struct {
	help        bool
	restart     bool
	userCommand string
	timeout     int
	yaml        bool
	quiet       bool
	verbose     bool
	file        string
	color       int
	debug       bool
	version     bool
	targets     []string
}

func usage1() string {
	return `Usage: gaze [options...] file(s)

Options(excerpt):
	-c  A command string.
	-r  Restart mode. Send SIGKILL to a ongoing process before invoking next.
	-t  Timeout(ms) Send SIGKILL to a ongoing process after this time.
	-h  Display help

Examples:
	gaze .
	gaze *.rb
	gaze main.go
	gaze -c make '**/*.c'
	gaze -c "eslint {{file}}" 'src/**/*.js'
	gaze -r server.py
	gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze
`
}

func usage2() string {
	return `Usage: gaze [options...] file(s)

Options:
	-c  A command string.
	-r  Restart mode. Send SIGKILL to a ongoing process before invoking next.
	-t  Timeout(ms) Send SIGKILL to a ongoing process after this time.
	-f  Specify a YAML configuration file.
	-v  Verbose mode.
	-q  Quiet mode.
	-y  Output the default configuration
	-h  Display help
	--color    Color(0:plain, 1:colorful)
	--version  Output version information

Examples:
	gaze .
	gaze *.rb
	gaze main.go
	gaze -c make '**/*.c'
	gaze -c "eslint {{file}}" 'src/**/*.js'
	gaze -r server.py
	gaze -t 1000 complicated.py

For more information: https://github.com/wtetsu/gaze
`
}
