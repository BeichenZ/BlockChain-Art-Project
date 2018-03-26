package RPiGPIO
//package main

import (
	bgpio "./gpio"
	"time"
)

const LEDOnMiliSecond = 500

const (
	NorthLED uint  = 5 
	SouthLED uint = 6
	EastLED  uint = 13
	WestLED  uint = 19
)
/*
//Test Program
func main() {
	for i:= 0; i<5;i++ {
		FlashLEDOnce(NorthLED)
	}

}
*/

func FlashLEDOnce(LEDName uint) {
	outputPin := bgpio.NewOutput(LEDName,false)
	outputPin.High()
	time.Sleep(LEDOnMiliSecond*time.Millisecond)
	outputPin.Low()
	time.Sleep(100*time.Millisecond)
}
