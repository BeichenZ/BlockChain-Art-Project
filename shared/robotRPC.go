package shared

import (
	"fmt"
	"time"
//	"encoding/json"
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
	*receivedMap = robotRPC.PiRobot.RMap
	//Testing
	//temp:= RandomMapGenerator()
	fmt.Println("RPC: RobotRPC:------> Sending map")
	fmt.Println(robotRPC.PiRobot.RMap)
	//*receivedMap = temp
	return nil
}
// This robot (server) adds the task given by the caller robot
func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *bool) error {
	robotRPC.PiRobot.ReceivedTasks = append(robotRPC.PiRobot.ReceivedTasks, *senderTask)

	fmt.Println("RPC:RobotRPC:---> ReceiveTASK")
	//data, _ :=json.MarshalIndent(senderTask, "", "")
	fmt.Println("RPC:SenderID", (*senderTask).SenderID, "=>", (*senderTask).DestPoint)

	var incommingMessage int
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTask.SendlogMessage, &incommingMessage)
	return nil
}

// TODO
func (robotRPC *RobotRPC) ReceiveTaskDecsionResponse(senderTaskDecision *TaskDescisionPayload, reply *ResponseForNeighbourPayload) error {
	var incommingMessage int
	fmt.Println("RPC:Receive task response from neighbour: ", senderTaskDecision.SenderAddr)
	robotRPC.PiRobot.ReceivedTasksResponse = append(robotRPC.PiRobot.ReceivedTasksResponse, *senderTaskDecision)
	fmt.Println("RPC: ",robotRPC.PiRobot.ReceivedTasksResponse)
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTaskDecision.SendlogMessage, &incommingMessage)
	return nil
}

func (robotRPC *RobotRPC) RegisterNeighbour(message *string, reply *string) error {
	// myNewNeighbour := Neighbour{Addr: *message}
	// robotRPC.PiRobot.RobotNeighbours = append(robotRPC.PiRobot.RobotNeighbours, myNewNeighbour)
	robotRPC.PiRobot.PossibleNeighbours.Add(*message) // hommie not used
	// robotRPC.PiRobot.PossibleNeighbours = append(robotRPC.PiRobot.PossibleNeighbours, *message)
	*reply = robotRPC.PiRobot.RobotIP
	//fmt.Println(*message)
	//fmt.Println("This is from ResgisterNeighbour")
	return nil
}

func (robotRPC *RobotRPC) NotifyNeighbours(p *Neighbour, ignore *bool) error {
	fmt.Printf("RPC:Adding neighbour %s to this robot %s \n", p.Addr, robotRPC.PiRobot.RobotIP )

	if robotRPC.PiRobot.RobotIP != p.Addr {
		robotRPC.PiRobot.RobotNeighbours.Lock()
		robotRPC.PiRobot.RobotNeighbours.rNeighbour[(*p).NID] = *p
		robotRPC.PiRobot.RobotNeighbours.Unlock()
		//robotRPC.PiRobot.RobotNeighbours = append(robotRPC.PiRobot.RobotNeighbours, *p)
	}
	return nil
}

func (r  *RobotStruct) WithinRadiusOfNetwork(p *FarNeighbourPayload) bool {
	r.RobotNeighbours.Lock()
	calleeNeighborCount := len(r.RobotNeighbours.rNeighbour)
	r.RobotNeighbours.Unlock()
	callerNeighborCount := len(p.ItsNeighbours)
	totalNodeCount := 2 + calleeNeighborCount + callerNeighborCount
	globalNodeArr := make([] Coordinate, totalNodeCount)
	globalNodeArr[0] = r.CurLocation
	globalNodeArr[1] = p.NeighbourCoordinate

	for i,callerNei := range p.ItsNeighbours {
		globalNodeArr[2+i] = callerNei.NeighbourCoordinate
	}

	j := 2 + callerNeighborCount
	r.RobotNeighbours.Lock()
	for _ , calleeNei := range r.RobotNeighbours.rNeighbour {
		globalNodeArr[j] = calleeNei.NeighbourCoordinate
		j++
	}
	r.RobotNeighbours.Unlock()

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

	return true
}

func (r *RobotStruct) RobotStateCommunicationAllowed(nid int) bool {
	var a bool
	r.RobotNeighbours.Lock()
	if _, ok := r.RobotNeighbours.rNeighbour[nid]; ok && (r.State.rState == BUSY) {
		r.RobotNeighbours.Unlock()
		a =true;
	}else {
		r.RobotNeighbours.Unlock()
		a = false
	}

	r.State.Lock()
	shit := r.State.rState == ROAM || r.State.rState == JOIN
	r.State.Unlock()



	return (a ||shit)
}


  // Server -> R2
// FN: Called by a robot to see if THIS robot is within its and its current neighbours CR
// Robot on this end will only return true if its in join or roam state (not if its in the busy state or if its in the roaming but the flag is of"
func (robotRPC *RobotRPC) ReceivePossibleNeighboursPayload(p *FarNeighbourPayload, responsePayload *ResponseForNeighbourPayload) error {
	fmt.Println("RPC: ReceivePossibleNeighboursPayload() robot Client that called this method and state (should be in roaming) ", p.NeighbourID, " ", p.State)
	var incommingMessage int
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", p.SendlogMessage, &incommingMessage)
	// TODO change this
	robotRPC.PiRobot.RobotNeighbours.Lock()
	for _, val :=range robotRPC.PiRobot.RobotNeighbours.rNeighbour{
		responsePayload.NeighboursNeighbourRobots = append(responsePayload.NeighboursNeighbourRobots, val)
	}
	robotRPC.PiRobot.RobotNeighbours.Unlock()
	// check on this later

	robotRPC.PiRobot.exchangeFlag.Lock()
	if !robotRPC.PiRobot.exchangeFlag.flag {
		robotRPC.PiRobot.exchangeFlag.Unlock()
		fmt.Println("RPC: FINISHED BUSY STATE. MUST WAIT UNTIL TIMER IS DONE TO TALK TO NEIGHBOR AGAIN")
		responsePayload.WithInComRadius = false
		return nil
	} else {
		robotRPC.PiRobot.exchangeFlag.Unlock()
	}
	fmt.Println("RPC: ReceivePossibleNeighboursPayload()  exchange flag ",robotRPC.PiRobot.exchangeFlag.flag, " robot within radius? ", robotRPC.PiRobot.WithinRadiusOfNetwork(p),
		"this robot state ", robotRPC.PiRobot.State.rState)

	//connection is formed only if the current robot is within CR and os either in ROAM or JOIN
	if robotRPC.PiRobot.WithinRadiusOfNetwork(p) && robotRPC.PiRobot.RobotStateCommunicationAllowed(p.NeighbourID){

		//newNeighbour.IsWithinCR  = true

		//fmt.Println("ReceivePossibleNeighboursPayload() Within the radius")
		//
		//fmt.Println()
		//fmt.Println("Join signal is sent.................................................")
		//fmt.Println("join sig received")
		//fmt.Println("neighbour IP is: ", newNeighbour.Addr)
		//fmt.Println("Waiting for the other robots to join")
		//fmt.Println()
		//

		//
		//fmt.Printf("Checking FirstTimeJoining: %x \n", robotRPC.PiRobot.joinInfo.firstTimeJoining)

		// This robot is sending ITS information
		rpcRobot := Neighbour{
			NID:                 robotRPC.PiRobot.RobotID,
			Addr:                robotRPC.PiRobot.RobotIP,
			NeighbourCoordinate: robotRPC.PiRobot.CurLocation,
			NMap:                robotRPC.PiRobot.RMap,
		}

		//responsePayload.NeighbourRobot = rpcRobot

		// put the robot itself into the NeighboursNeighbourRobots
		responsePayload.NeighboursNeighbourRobots = append(responsePayload.NeighboursNeighbourRobots, rpcRobot)

		//responsePayload.NeighboursNeighbourRobots = robotRPC.PiRobot.RobotNeighbours
		responsePayload.NeighbourState = robotRPC.PiRobot.State.rState

		// robotRPC.PiRobot.RobotNeighbours[newNeighbour.NID] = newNeighbour


		robotRPC.PiRobot.State.Lock()
		if robotRPC.PiRobot.State.rState == JOIN {
			robotRPC.PiRobot.State.Unlock()
			responsePayload.RemainingTime = time.Now().Sub(robotRPC.PiRobot.joinInfo.joiningTime)
		} else {
			robotRPC.PiRobot.State.Unlock()
		}
		responsePayload.WithInComRadius = true

		//// This robot (server) will add the client and its neighbours to itself
		//SaveNeighbour(robotRPC.PiRobot, p.ItsNeighbours)

	}else{
		//skip the request client
		responsePayload.WithInComRadius = false
	}

	return nil
}
