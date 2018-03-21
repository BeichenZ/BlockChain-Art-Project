package shared

import "fmt"

type RobotRPC struct {
	PiRobot *RobotStruct
}

func (robotRPC *RobotRPC) ReceiveMap(senderMap *Map, reply *int) error {
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *int) error {
	fmt.Println(senderTask.SendlogMessage)
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTask.SendlogMessage, TaskPayload{})
	return nil
}

func (robotRPC *RobotRPC) RegisterNeighbour(message *string, reply *string) error {
	// myNewNeighbour := Neighbour{Addr: *message}
	// robotRPC.PiRobot.RobotNeighbours = append(robotRPC.PiRobot.RobotNeighbours, myNewNeighbour)
	robotRPC.PiRobot.PossibleNeighbours = append(robotRPC.PiRobot.PossibleNeighbours, *message)
	*reply = robotRPC.PiRobot.RobotIP
	fmt.Println(*message)
	return nil
}
