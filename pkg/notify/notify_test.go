/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package notify

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/wtetsu/gaze/pkg/logger"
)

func TestBasic(t *testing.T) {
	logger.Level(logger.VERBOSE)

	rb := createTempFile("*.rb", `puts "Hello from Ruby`)
	py := createTempFile("*.py", `print("Hello from Python")`)

	pattens := []string{rb, py}

	notify, err := New(pattens)
	if err != nil {
		t.Fatal()
	}

	count := 0
	go func() {
		for {
			select {
			case _, ok := <-notify.Events:
				count++
				if !ok {
					continue
				}
			case err, ok := <-notify.Errors:
				count++
				if !ok {
					continue
				}
				log.Println("error:", err)
			}
		}
	}()

	for i := 0; i < 100; i++ {
		touch(py)
		touch(rb)
		time.Sleep(10)
		if count >= 2 {
			break
		}
	}
	if count == 0 {
		t.Fatal()
	}

	notify.Close()
}

func createTempFile(pattern string, content string) string {
	file, err := ioutil.TempFile("", pattern)
	if err != nil {
		return ""
	}
	file.WriteString(content)
	file.Close()

	return file.Name()
}

func touch(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	file.WriteString(" ")
	file.Close()
}
