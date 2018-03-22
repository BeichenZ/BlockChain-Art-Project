package shared

import "fmt"

type RobotRPC struct {
	PiRobot *RobotStruct
}

type FarNeighbourPayload struct {
	NeighbourID         int
	NeighbourCoordinate PointStruct
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
	robotRPC.PiRobot.PossibleNeighbours.Add(*message)
	// robotRPC.PiRobot.PossibleNeighbours = append(robotRPC.PiRobot.PossibleNeighbours, *message)
	*reply = robotRPC.PiRobot.RobotIP
	fmt.Println(*message)
	return nil
}

// This funciton is periodically called to detemine the distance between two neighbours
func (robotRPC *RobotRPC) ReceivePossibleNeighboursPayload(p *FarNeighbourPayload, reply *string) error {
	// Calculate distance here
	fmt.Println("Getting neighbour info")
	fmt.Println(p.NeighbourID)
	// distance := 0
	// if distance < 1 {
	// 	robotRPC.PiRobot.JoiningSig <- p.NeighbourID
	// }
	return nil
}
