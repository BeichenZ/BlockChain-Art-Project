package shared

import "fmt"

type RobotRPC struct {
	PiRobot *RobotStruct
}

type FarNeighbourPayload struct {
	NeighbourID         int
	NeighbourIPAddr     string
	NeighbourCoordinate Coordinate
	NeighbourMap        Map
	SendlogMessage      []byte
}

func (robotRPC *RobotRPC) ReceiveMap(senderMap *Map, reply *int) error {
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *int) error {
	fmt.Println(senderTask.SendlogMessage)
	var incommingMessage int
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTask.SendlogMessage, &incommingMessage)
	return nil
}

// TODO
func (robotRPC *RobotRPC) ReceiveTaskDecsionResponse(senderTaskDecision *TaskDescisionPayload, reply *int) error {
	var incommingMessage int
	fmt.Println("Receive task response from neighbour: ", senderTaskDecision.SenderAddr)
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTaskDecision.SendlogMessage, &incommingMessage)
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
	var incommingMessage int
	fmt.Println("receive info from neighbour: ", p.NeighbourID)
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", p.SendlogMessage, &incommingMessage)
	// TODO change this
	newNeighbour := Neighbour{
		NID:                 p.NeighbourID,
		Addr:                p.NeighbourIPAddr,
		NeighbourCoordinate: p.NeighbourCoordinate,
		NMap:                p.NeighbourMap,
	}
	distance := 0
	if distance < 1 {
		robotRPC.PiRobot.JoiningSig <- newNeighbour
	}
	return nil
}
