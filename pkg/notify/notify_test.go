/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package notify

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/time"
)

func TestUtilFunctions(t *testing.T) {
	logger.Level(logger.VERBOSE)

	tmpDir := createTempDir()

	if tmpDir == "" {
		t.Fatal("Temp files error")
	}

	os.MkdirAll(tmpDir+"/dir0", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2a/dir3a", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2a/dir3b", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2a/dir3c", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2b/dir3a", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2b/dir3b", os.ModePerm)
	os.MkdirAll(tmpDir+"/dir1/dir2b/dir3c", os.ModePerm)

	createTempFileWithDir(tmpDir+"/dir1/dir2b/dir3b", "*.tmp", `puts "Hello from Ruby`)

	actual1 := findDirs([]string{tmpDir + "/*"})
	sort.Strings(actual1)

	expected1 := []string{
		tmpDir,
		tmpDir + "/dir0",
		tmpDir + "/dir1",
	}

	for i := 0; i < len(expected1); i++ {

		if filepath.Clean(actual1[i]) != filepath.Clean(expected1[i]) {
			t.Fatalf("%s != %s", actual1[i], expected1[i])
		}
	}

	actual2 := findDirs([]string{tmpDir + "/**"})
	sort.Strings(actual2)

	expected2 := []string{
		tmpDir,
		tmpDir + "/dir0",
		tmpDir + "/dir1",
		tmpDir + "/dir1/dir2a",
		tmpDir + "/dir1/dir2a/dir3a",
		tmpDir + "/dir1/dir2a/dir3b",
		tmpDir + "/dir1/dir2a/dir3c",
		tmpDir + "/dir1/dir2b",
		tmpDir + "/dir1/dir2b/dir3a",
		tmpDir + "/dir1/dir2b/dir3b",
		tmpDir + "/dir1/dir2b/dir3c",
	}

	for i := 0; i < len(expected2); i++ {

		if filepath.Clean(actual2[i]) != filepath.Clean(expected2[i]) {
			t.Fatalf("%s != %s", actual2[i], expected2[i])
		}
	}
}

func TestFindRealDirectory(t *testing.T) {
	tmpDir := createTempDir()

	if tmpDir == "" {
		t.Fatal("Temp files error")
	}

	os.MkdirAll(tmpDir+"/dir0/dir1/dir2/dir3", os.ModePerm)

	var r string

	r = findRealDirectory(tmpDir + "/dir0/dir1/dir2/dir3")
	if r != filepath.Clean(tmpDir+"/dir0/dir1/dir2/dir3") {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory(tmpDir + "/dir0/dir1/dir2/**")
	if r != filepath.Clean(tmpDir+"/dir0/dir1/dir2") {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory(tmpDir + "/dir0/dir1/")
	if r != filepath.Clean(tmpDir+"/dir0/dir1") {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory(tmpDir + "/dir0/dir1/**/dir3")
	if r != filepath.Clean(tmpDir+"/dir0/dir1") {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory(tmpDir + "/?ir?/dir1/dir2/dir3")
	if r != filepath.Clean(tmpDir) {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory(tmpDir + "/dir0/dir1/\\*\\?\\[\\]/dir3")
	if r != filepath.Clean(tmpDir+"/dir0/dir1") {
		t.Fatal("Unexpected result:" + r)
	}

	r = findRealDirectory("invalid/path/")
	if r != "" {
		t.Fatal("Unexpected result:" + r)
	}
}

func TestUpdate(t *testing.T) {
	logger.Level(logger.VERBOSE)

	rb := createTempFile("*.rb", `puts "Hello from Ruby`)
	py := createTempFile("*.py", `print("Hello from Python")`)

	if rb == "" || py == "" {
		t.Fatal("Temp files error")
	}

	pattens := []string{filepath.Dir(rb) + "/*.rb", filepath.Dir(rb) + "/*.py"}

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
		if count >= 2 {
			break
		}
		time.Sleep(20)
	}
	if count < 2 {
		t.Fatalf("count:%d", count)
	}

	notify.Close()
	notify.Close()
}

func TestCreateAndMove(t *testing.T) {
	logger.Level(logger.VERBOSE)

	tmpDir := createTempDir()

	if tmpDir == "" {
		t.Fatal("Temp files error")
	}

	notify, err := New([]string{tmpDir})
	notify.regardRenameAsModPeriod = 10000
	notify.detectCreate = true
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
		rb := createTempFileWithDir(tmpDir, "*.tmp", `puts "Hello from Ruby`)
		os.Rename(rb, rb+".rb")
		py := createTempFileWithDir(tmpDir, "*.tmp", `print("Hello from Python")`)
		os.Rename(py, py+".py")

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

func TestDelete(t *testing.T) {
	logger.Level(logger.VERBOSE)

	rb1 := createTempFile("*.rb", `puts "Hello from Ruby`)
	rb2 := createTempFile("*.rb", `puts "Hello from Ruby`)
	py1 := createTempFile("*.py", `print("Hello from Python")`)
	py2 := createTempFile("*.py", `print("Hello from Python")`)

	if rb1 == "" || rb2 == "" || py1 == "" || py2 == "" {
		t.Fatal("Temp files error")
	}

	pattens := []string{
		filepath.Dir(rb1) + "/*.rb",
		filepath.Dir(rb2) + "/*.rb",
		filepath.Dir(py1) + "/*.py",
		filepath.Dir(py2) + "/*.py",
	}

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

	os.Remove(rb1)
	os.Remove(rb2)
	os.Remove(py1)
	os.Remove(py2)

	time.Sleep(20)

	if count != 0 {
		t.Fatalf("count:%d", count)
	}

	notify.Close()
	notify.Close()
}

func TestQueue(t *testing.T) {
	logger.Level(logger.VERBOSE)

	rb := createTempFile("*.rb", `puts "Hello from Ruby`)
	py := createTempFile("*.py", `print("Hello from Python")`)

	if rb == "" || py == "" {
		t.Fatal("Temp files error")
	}

	rbCommand := fmt.Sprintf(`ruby "%s"`, rb)
	pyCommand := fmt.Sprintf(`python "%s"`, py)

	pattens := []string{filepath.Dir(rb) + "/*.rb", filepath.Dir(rb) + "/*.py"}

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

	notify.Requeue(Event{rbCommand, 3})
	notify.Requeue(Event{pyCommand, 4})
	notify.Requeue(Event{rbCommand, 5})
	notify.Requeue(Event{pyCommand, 6})
	for i := 0; i < 50; i++ {
		// touch(py)
		// touch(rb)
		if count >= 2 {
			break
		}
		time.Sleep(20)
	}
	if count < 2 {
		t.Fatalf("count:%d", count)
	}

	notify.Close()
	notify.Close()
}

func createTempDir() string {
	dirpath, err := ioutil.TempDir("", "_gaze")
	if err != nil {
		return ""
	}
	return dirpath
}

func createTempFile(pattern string, content string) string {
	dirpath := createTempDir()
	return createTempFileWithDir(dirpath, pattern, content)
}

func createTempFileWithDir(dirpath string, pattern string, content string) string {
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
