package main //package RPiGPIO

import (
	bgpio "./gpio"
	"fmt"
	"time"
	"os" )
const ButtonWaitMiliSecond = 300

//Used to distinguish between different button pressed
type ButtonType int
const (
	FrontObstacleButton ButtonType = 1+iota
	FrontEmptyButton
	LeftObstacleButton
	RightObstacleButton

)
const (
	FrontObstacleButton_Pin uint = 5
	FrontEmptyButton_Pin uint = 6
	LeftObstacleButton_Pin uint = 13
	RightObstacleButton_Pin uint =19 )

func main() {
	//Test Codes for button
	recChannel := make(chan int)
	MonitorButtonAtPin(LeftObstacleButton,recChannel)	
	temp := <- recChannel
	fmt.Printf("Button is pressed with Enum: %d\n",temp)
	
}


//Note:Function Assumes Monitored Pin has been pulled down
func MonitorButtonAtPin(button ButtonType,outputChan chan int ){
	var pinNum uint
	switch button {
		case FrontObstacleButton:
			pinNum = FrontObstacleButton_Pin
		case FrontEmptyButton:
			pinNum = FrontEmptyButton_Pin
		case LeftObstacleButton:
			pinNum = LeftObstacleButton_Pin
		case RightObstacleButton:
			pinNum = RightObstacleButton_Pin
		default:
		     fmt.Println("Unrecognized Button Pressed with ButtonType Enum:",button)
		     os.Exit(1)

	}
	buttonPin := bgpio.NewInput(pinNum)
	go func (){
		Loop:
			for{
                        	//Button is pressed when input is high(skip pull down)
				if value,err := buttonPin.Read();value==1{
					fmt.Printf("Button at Pin%d is pressed\n",pinNum)
					if err != nil {fmt.Println(err)}
					time.Sleep(ButtonWaitMiliSecond*time.Millisecond)
					outputChan <- int(button)
         				        break Loop //Terminate after button is pressed once								  }		
                		}
			  }
	}()
}

