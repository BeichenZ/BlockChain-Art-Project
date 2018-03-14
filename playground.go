package main

import (
	"./shared"
	"fmt"
	"bufio"
	"os"
	"time"
)

func main()  {

	robot :=shared.InitRobot(123, shared.Map{})
	go robot.Explore()

	go func() {
		for{

			fmt.Println("YOYOYO")
			robot.SendFreeSpaceSig()
			time.Sleep(1000*time.Millisecond)
		}
	}()


	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}