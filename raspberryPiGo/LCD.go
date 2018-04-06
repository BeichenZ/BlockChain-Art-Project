//package main
package RPiGPIO
import (
	"fmt"
	hd44780 "./go-hd44780"
	"time"
)

func main(){
	lcd := hd44780.NewGPIO4bit()
	fmt.Println("Start Display Set up")
	if err := lcd.Open();err != nil {
		panic("Cannot OPen lcd:"+err.Error())
	}
	lcd.DisplayLines("line1\nline2")
	//time.Sleep(2000*time.Millisecond)
	//lcd.DisplayLines("State:JOin")
	//lcd.DisplayLines("\nSecond LIne")
	//time.Sleep(2000*time.Millisecond)
	//lcd.Close()
}

type LCDScreen struct {
	name string
}
func (lcdObj* LCDScreen) LCD_DisplayString(input string){
	//TODO: init and close of LCD should be done in the robot level but not every time!
	//TODO: Refactor the code
	lcd := hd44780.NewGPIO4bit()
        if err := lcd.Open();err != nil {
                fmt.Println("Cannot OPen lcd,Message is lost:",input)
        }
        lcd.DisplayLines(input)
        time.Sleep(2*time.Second)
        lcd.Close()
}
