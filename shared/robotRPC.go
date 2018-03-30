package shared

import (
	"fmt"
	"time"
	"encoding/json"
)

type RobotRPC struct {
	PiRobot *RobotStruct
}

type FarNeighbourPayload struct {
	NeighbourID         int
	NeighbourIPAddr     string
	NeighbourCoordinate Coordinate
	NeighbourMap        Map
	SendlogMessage      []byte
	State               RobotState
	ItsNeighbours       []Neighbour
}

type ResponseForNeighbourPayload struct {
	WithInComRadius bool
	RemainingTime  time.Duration
	NeighbourRobot Neighbour
	NeighbourState RobotState
	NeighboursNeighbourRobots []Neighbour

}

func (robotRPC *RobotRPC) ReceiveMap(ignore bool, receivedMap *Map) error {
	//Productio code
	//*receivedMap = robotRPC.PiRobot.RMap
	//Testing
	temp:= RandomMapGenerator()
	fmt.Println("RobotRPC:------> Sending random map")
	fmt.Println(temp)
	*receivedMap = temp
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *int) error {
	robotRPC.PiRobot.ReceivedTasks = append(robotRPC.PiRobot.ReceivedTasks, *senderTask)

	fmt.Println("RobotRPC: ReceiveTASK--->")
	data, _ :=json.MarshalIndent(senderTask, "", "")
	fmt.Println(string(data))

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
	//fmt.Println(*message)
	//fmt.Println("This is from ResgisterNeighbour")
	return nil
}

func (robotRPC *RobotRPC) NotifyNeighbours(p *Neighbour, ignore *bool) error {
	fmt.Printf("Adding neighbour %s to this robot %s \n", p.Addr, robotRPC.PiRobot.RobotIP )

	if robotRPC.PiRobot.RobotIP != p.Addr {
		robotRPC.PiRobot.RobotNeighbours[(*p).NID] = *p
		//robotRPC.PiRobot.RobotNeighbours = append(robotRPC.PiRobot.RobotNeighbours, *p)
	}
	return nil
}

func (r  *RobotStruct) WithinRadiusOfNetwork(p *FarNeighbourPayload) bool {

	calleeNeighborCount := len(r.RobotNeighbours)
	callerNeighborCount := len(p.ItsNeighbours)
	totalNodeCount := 2 + calleeNeighborCount + callerNeighborCount
	globalNodeArr := make([] Coordinate, totalNodeCount)
	globalNodeArr[0] = r.CurLocation
	globalNodeArr[1] = p.NeighbourCoordinate
	for i,callerNei := range p.ItsNeighbours {
		globalNodeArr[2+i] = callerNei.NeighbourCoordinate
	}

	j := 2 + callerNeighborCount
	for _ , calleeNei := range r.RobotNeighbours{
		globalNodeArr[j] = calleeNei.NeighbourCoordinate
		j++
	}

	for  i := 0; i < len(globalNodeArr); i++ {
		ithNodeCoordinate := globalNodeArr[i]
		for j := i+1; j < len(globalNodeArr); j++ {
			jthNodeCoordinate := globalNodeArr[j]

			dist := DistBtwnTwoPoints(jthNodeCoordinate, ithNodeCoordinate)

			if dist > EXRADIUS {
				return false
			}
		}
	}

	return true && (r.State == ROAM || r.State == JOIN)
}


  // Server -> R2
// This funciton is periodically called to detemine the distance between two neighbours
func (robotRPC *RobotRPC) ReceivePossibleNeighboursPayload(p *FarNeighbourPayload, responsePayload *ResponseForNeighbourPayload) error {
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
		IsWithinCR:			 false,
	}

	for _, val :=range robotRPC.PiRobot.RobotNeighbours{
		responsePayload.NeighboursNeighbourRobots = append(responsePayload.NeighboursNeighbourRobots, val)
	}

	if !robotRPC.PiRobot.exchangeFlag.flag {
		responsePayload.WithInComRadius = false
		return nil
	}

	//connection is formed only if the current robot is within CR and os either in ROAM or JOIN
	if robotRPC.PiRobot.WithinRadiusOfNetwork(p){



		newNeighbour.IsWithinCR  = true

		fmt.Println("ReceivePossibleNeighboursPayload() Within the radius")


		fmt.Println("Join signal is sent.................................................")
		fmt.Println("join sig received")
		fmt.Println("neighbour IP is: ", newNeighbour.Addr)
		fmt.Println("Waiting for the other robots to join")

		responsePayload.WithInComRadius = true

		fmt.Printf("Checking FirstTimeJoining: %x \n", robotRPC.PiRobot.joinInfo.firstTimeJoining)


		rpcRobot := Neighbour{
			NID:                 robotRPC.PiRobot.RobotID,
			Addr:                robotRPC.PiRobot.RobotIP,
			NeighbourCoordinate: robotRPC.PiRobot.CurLocation,
			NMap:                robotRPC.PiRobot.RMap,
			IsWithinCR:			 true,
		}

		responsePayload.NeighbourRobot = rpcRobot

		// put the robot itself into the NeighboursNeighbourRobots
		responsePayload.NeighboursNeighbourRobots = append(responsePayload.NeighboursNeighbourRobots, rpcRobot)

		//responsePayload.NeighboursNeighbourRobots = robotRPC.PiRobot.RobotNeighbours
		responsePayload.NeighbourState = robotRPC.PiRobot.State

		robotRPC.PiRobot.RobotNeighbours[newNeighbour.NID] = newNeighbour

		if robotRPC.PiRobot.State == JOIN {
			responsePayload.RemainingTime = time.Now().Sub(robotRPC.PiRobot.joinInfo.joiningTime)
			fmt.Println("Remaining Time is ", responsePayload.RemainingTime)
		}

	}else{
		//skip the request client
		responsePayload.WithInComRadius = false
	}

	return nil
}
