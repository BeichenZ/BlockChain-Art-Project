package main

import (
	shared "../shared"
	"github.com/stianeikeland/go-rpio"
	"fmt"
	"os"
	"time"
)

var (
	pinNorth = rpio.Pin(10)
	pinSouth = rpio.Pin(11)
	pinEast = rpio.Pin(25)
	pinWest = rpio.Pin(8)

	pinFreeSpaceButton = rpio.Pin(5)
	pinWallButton = rpio.Pin(6)
	pinRightBumperButton = rpio.Pin(12)


)
func main () {
	// Pin setup
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pins to output mode
	pinNorth.Output()
	pinSouth.Output()
	pinEast.Output()
	pinWest.Output()

	path := [...]shared.PointStruct{shared.SOUTH, shared.NORTH, shared.WEST, shared.EAST, shared.SOUTH, shared.NORTH, shared.WEST, shared.EAST}

	for _,dir := range path{
		fmt.Println(1)
		switch dir {
			case shared.NORTH: {
				pinNorth.High()
				time.Sleep(2*time.Second)
				pinNorth.Low()
				time.Sleep(1*time.Second)
				break;
			}
			case shared.SOUTH: {
				pinSouth.High()
				time.Sleep(2*time.Second)
				pinSouth.Low()
				time.Sleep(1*time.Second)
				break;
			}
			case shared.EAST: {
				pinEast.High()
				time.Sleep(2*time.Second)
				pinEast.Low()
				time.Sleep(1*time.Second)
				break;
			}
			case shared.WEST: {
				pinWest.High()
				time.Sleep(2*time.Second)
				pinWest.Low()
				time.Sleep(1*time.Second)
				break;
			}

		}
	}

}

