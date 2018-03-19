package main

import (
	bgpio "github.com/brian-armstrong/gpio"
	"fmt"
)

func main() {
	watcher := bgpio.NewWatcher()
	watcher.AddPin(10)
	defer watcher.Close()

	go func() {
		for {
			pin, value := watcher.Watch()
			fmt.Printf("read %d from gpio %d\n", value, pin)
		}
	}()

	for {
		_ = 2;
	}
}

