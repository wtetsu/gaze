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

	"github.com/wtetsu/gaze/pkg/app"
	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/logger"
)

func main() {
	help := flag.Bool("h", false, "help")
	// parallel := flag.Bool("p", false, "bool flag")
	// recursion := flag.Bool("r", false, "bool flag")
	userCommand := flag.String("c", "", "command")
	timeout := flag.Int("t", 0, "int flag")
	yaml := flag.Bool("y", false, "Show default config")
	quiet := flag.Bool("q", false, "")
	verbose := flag.Bool("v", false, "")
	file := flag.String("f", "", "")
	color := flag.Int("color", 1, "")

	flag.Parse()

	if *yaml {
		fmt.Println(config.Default())
		return
	}

	if *help || len(flag.Args()) == 0 {
		fmt.Println(usage())
		return
	}

	if *color == 0 {
		logger.Plain()
	} else {
		logger.Colorful()
	}
	if *quiet {
		logger.Level(logger.QUIET)
	}
	if *verbose {
		logger.Level(logger.VERBOSE)
	}

	err := app.Start(flag.Args(), *userCommand, *file, *timeout)
	if err != nil {
		logger.ErrorObject(err)
		os.Exit(1)
	}
}

func usage() string {
	return `
	Usage: gaze [files...] [options...]
	
	Options:
		-c  Command.
		-q  Quiet.
		-r  Recursive.
		-p  Parallel.
		f  File.
	-y  Show default configuration. Save as ./.gaze.yml or ~/.gaze.yml and edit it.
	`
}
