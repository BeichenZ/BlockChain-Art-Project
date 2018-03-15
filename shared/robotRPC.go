package shared

import "fmt"

type RobotRPC struct {
	PiRobot *RobotStruct
}

func (robotRPC *RobotRPC) ReceiveMap(senderMap *Map, reply *int) error {
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask string, reply *int) error {
	fmt.Println(senderTask)
	return nil
}
