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
	// quiet := flag.Bool("q", false, "bool flag")
	// parallel := flag.Bool("p", false, "bool flag")
	// recursion := flag.Bool("r", false, "bool flag")
	userCommand := flag.String("c", "", "command")
	yaml := flag.Bool("y", false, "Show default config")
	flag.Parse()

	if *help {
		fmt.Println(usage)
		return
	}

	if *yaml {
		fmt.Println(config.Default())
		return
	}

	err := app.Start(flag.Args(), *userCommand)
	if err != nil {
		logger.Debug(err)
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
