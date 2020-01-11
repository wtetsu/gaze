package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
)

type Command struct {
	shell    string
	fileMap  map[string]string
	disposed bool
}

// New returns a new command
func New(shell string) *Command {
	return &Command{
		shell:    shell,
		fileMap:  make(map[string]string),
		disposed: false,
	}
}

func (c *Command) PrepareScript(commandString string) string {
	if c.disposed {
		return ""
	}
	existingFilePath, found := c.fileMap[commandString]

	if found && fs.IsFile(existingFilePath) {
		return existingFilePath
	}

	newFilePath, err := ioutil.TempFile("", "*.gaze.cmd")
	if err != nil {
		return ""
	}

	if c.shell == "cmd" {
		newFilePath.WriteString("@" + commandString)
	} else {
		newFilePath.WriteString(commandString)
	}
	err = newFilePath.Close()
	if err != nil {
		return ""
	}

	c.fileMap[commandString] = newFilePath.Name()

	return newFilePath.Name()
}

// Dispose disposes all the internal resources
func (c *Command) Dispose() {
	if c.disposed {
		return
	}
	for _, filePath := range c.fileMap {
		fmt.Println(filePath)
		err := os.Remove(filePath)
		if err != nil {
			logger.DebugObject(err)
		}
	}
	c.disposed = true
}

func (c *Command) Shell() string {
	return c.shell
}
func (c *Command) Disposed() bool {
	return c.disposed
}
