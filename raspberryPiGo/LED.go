//package RPiGPIO
package RPiGPIO

import (
	bgpio "./gpio"
	"time"
)

const LEDOnMiliSecond = 500

const (
	NorthLED uint  = 4 //green
	SouthLED uint = 17 //red
	EastLED  uint = 27 //yellow
	WestLED  uint = 22 //blue
)

//Test Program
func main() {
	for i:= 0; i<5;i++ {
		FlashLEDOnce(NorthLED)
		FlashLEDOnce(SouthLED)
		FlashLEDOnce(EastLED)
		FlashLEDOnce(WestLED)
	}

}


func FlashLEDOnce(LEDName uint) {
	outputPin := bgpio.NewOutput(LEDName,false)
	outputPin.High()
	time.Sleep(LEDOnMiliSecond*time.Millisecond)
	outputPin.Low()
	time.Sleep(100*time.Millisecond)
}
