/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package gazer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/time"
)

func TestBasic(t *testing.T) {
	py1 := createTempFile("*.py", `print("hello!!!")`)
	txt1 := createTempFile("*.txt", `print("hello!!!")`)

	gazer := New([]string{py1, txt1})
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.InitConfig([]string{".gaze.yml", ".gaze.yaml"})
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 0, false)
	time.Sleep(100)

	if gazer.Counter() != 0 {
		t.Fatal()
	}

	touch(py1)
	touch(txt1)
	time.Sleep(100)
	if gazer.Counter() != 1 {
		t.Fatal()
	}
}

func TestRestart(t *testing.T) {
	content := `
import time

print("start")
# time.sleep(1)
print("end")
`

	py1 := createTempFile("*.py", content)

	gazer := New([]string{py1})
	if gazer == nil {
		t.Fatal()
	}
	defer gazer.Close()

	c, err := config.InitConfig([]string{".gaze.yml", ".gaze.yaml"})
	if err != nil {
		t.Fatal()
	}
	go gazer.Run(c, 0, true)

	time.Sleep(100)

	if gazer.Counter() != 0 {
		t.Fatal()
	}

	touch(py1)
	time.Sleep(100)
	touch(py1)
	time.Sleep(100)
	touch(py1)
	time.Sleep(100)

	if gazer.Counter() != 3 {
		t.Fatal()
	}
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
	file.WriteString("#\n")
	file.Close()
}
