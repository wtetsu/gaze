package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	// "github.com/wtetsu/gaze/pkg/gaze"
)

func main() {
	log.Println("Start!")

	// ExampleNewWatcher()

	log.Println("End!")
}

func ExampleNewWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./aaa.txt")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
