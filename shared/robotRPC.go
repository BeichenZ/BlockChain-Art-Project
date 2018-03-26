package shared

import (
	"fmt"
	"time"
	"golang.org/x/tools/go/gcimporter15/testdata"
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

}

func (robotRPC *RobotRPC) ReceiveMap(ignore bool, receivedMap *Map) error {
	*receivedMap = robotRPC.PiRobot.RMap
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *int) error {
	robotRPC.PiRobot.ReceivedTasks = append(robotRPC.PiRobot.ReceivedTasks, *senderTask)
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
	//fmt.Println(*message)
	//fmt.Println("This is from ResgisterNeighbour")
	return nil
}
  // Server -> R3
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
	}
	distance := 0
	
	if distance < 1 && p.State == ROAM {
		fmt.Println("Join signal is sent.................................................")

		go func () {
			robotRPC.PiRobot.JoiningSig <- newNeighbour
		}()
		responsePayload.WithInComRadius = true

		fmt.Printf("Checking FirstTimeJoining: %x \n", robotRPC.PiRobot.joinInfo.firstTimeJoining)


		if robotRPC.PiRobot.joinInfo.firstTimeJoining {

			robotRPC.PiRobot.joinInfo.firstTimeJoining = false

			fmt.Println("Starting Time....................................")
			robotRPC.PiRobot.joinInfo.joiningTime = time.Now()
			fmt.Println(robotRPC.PiRobot.joinInfo.joiningTime)

			go func() {
				for {
					if time.Now().Sub(robotRPC.PiRobot.joinInfo.joiningTime) >= (TIMETOJOIN) {
						fmt.Println("Timer has ended. Going to the BUSY state..............")
						robotRPC.PiRobot.joinInfo.firstTimeJoining = true
						robotRPC.PiRobot.BusySig <- true
						break
					}
				}
			}()
		} else {

			responsePayload.RemainingTime = time.Now().Sub(robotRPC.PiRobot.joinInfo.joiningTime)

			fmt.Println("Remaining Time is ", responsePayload.RemainingTime)
			responsePayload.NeighbourRobot = Neighbour{
				NID:                 robotRPC.PiRobot.RobotID,
				Addr:                robotRPC.PiRobot.RobotIP,
				NeighbourCoordinate: robotRPC.PiRobot.CurLocation,
				NMap:                robotRPC.PiRobot.RMap,
			}
			responsePayload.NeighbourState = robotRPC.PiRobot.State

		}

	}

	return nil
}
