package main

import (
	"fmt"
	"github.com/wtetsu/gaze/pkg/app"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Usage: gaze command [file1] [file2] [...]")
		return
	}

	app.Start(os.Args[1:])
}
