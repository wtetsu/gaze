package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/wtetsu/gaze/pkg/gazer"
)

func main() {
	fmt.Println("Starting double Close() race condition test...")

	// Create Gazer instance
	g, err := gazer.New([]string{"."}, 100)
	if err != nil {
		panic(err)
	}

	// Call Close() simultaneously from multiple goroutines
	var wg sync.WaitGroup

	// 10 goroutines call Close() simultaneously
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d: Starting Close() call\n", id)
			g.Close()
			fmt.Printf("Goroutine %d: Finished Close() call\n", id)
		}(i)
	}

	// Wait a bit before execution
	time.Sleep(10 * time.Millisecond)

	// Also call Close() sequentially
	g.Close()
	g.Close()

	wg.Wait()
	fmt.Println("✅ Double Close() test completed! No panic!")
}
