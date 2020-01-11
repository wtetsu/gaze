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
	// timeout := flag.Int("t", 0, "int flag")
	// parallel := flag.Bool("p", false, "bool flag")
	// recursion := flag.Bool("r", false, "bool flag")
	userCommand := flag.String("c", "", "command")
	yaml := flag.Bool("y", false, "Show default config")
	quiet := flag.Bool("q", false, "")
	verbose := flag.Bool("v", false, "")
	flag.Parse()

	if *yaml {
		fmt.Println(config.Default())
		return
	}

	if *help || len(flag.Args()) == 0 {
		fmt.Println(usage)
		return
	}

	if *quiet {
		logger.Level(logger.QUIET)
	}
	if *verbose {
		logger.Level(logger.VERBOSE)
	}

	err := app.Start(flag.Args(), *userCommand)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

var usage = `
Usage: gaze [files...] [options...]

Options:
  -c  Command.
  -q  Quiet.
  -r  Recursive.
  -p  Parallel.
  -f  File.
  -y  Show default configuration. Save as ./.gaze.yml or ~/.gaze.yml and edit it.
`
