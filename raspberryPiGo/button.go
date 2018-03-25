package main

import (
	bgpio "./gpio"
	"fmt"
)

func main() {
	pin := bgpio.NewInput(19)
	pin.High()
	watcher := bgpio.NewWatcher()
	watcher.AddPin(19)
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

