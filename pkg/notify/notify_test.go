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

	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

func TestBasic(t *testing.T) {
	logger.Level(logger.VERBOSE)

	rb := createTempFile("*.rb", `puts "Hello from Ruby`)
	py := createTempFile("*.py", `print("Hello from Python")`)

	if rb == "" || py == "" {
		t.Fatal("Temp files error")
	}

	pattens := []string{rb, py}

	notify, err := New(pattens)
	if err != nil {
		t.Fatal()
	}

	notify.PendingPeriod(10)

	count := 0
	go func() {
		for {
			select {
			case _, ok := <-notify.Events:
				if !ok {
					continue
				}
				count++

			case err, ok := <-notify.Errors:
				if !ok {
					continue
				}
				log.Println("error:", err)
				count++
			}
		}
	}()

	for i := 0; i < 50; i++ {
		touch(py)
		touch(rb)
		if count >= 4 {
			break
		}
		time.Sleep(20)
	}
	if count < 4 {
		t.Fatalf("count:%d", count)
	}

	notify.Close()
	notify.Close()
}

func createTempFile(pattern string, content string) string {
	dirpath, err := ioutil.TempDir("", "_gaze")

	if err != nil {
		return ""
	}

	file, err := ioutil.TempFile(dirpath, pattern)
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
