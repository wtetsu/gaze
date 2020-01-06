package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wtetsu/gaze/pkg/app"
)

func main() {
	help := flag.Bool("h", false, "help")
	// timeout := flag.Int("t", 0, "int flag")
	// quiet := flag.Bool("q", false, "bool flag")
	// parallel := flag.Bool("p", false, "bool flag")
	// recursion := flag.Bool("r", false, "bool flag")
	userCommand := flag.String("c", "", "command")
	flag.Parse()

	if *help {
		fmt.Println(usage)
		return
	}

	err := app.Start(flag.Args(), *userCommand)
	if err != nil {
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
  -f  Filter.
`
