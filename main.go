package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"./shared"
)

// TODO: Include golang GPIO

func main() {
	/// Need to change to different ip address. May to use a different library due to ad-hoc
	IPAddr := os.Args[1]
	resolvedIPAddress, error := net.ResolveTCPAddr("tcp", IPAddr)
	if error != nil {
		log.Fatal("Unable to resolve IP Address", error)
	}

	listener, error := net.ListenTCP("tcp", resolvedIPAddress)

	if error != nil {
		log.Fatal("Unable to create a listner", error)
	}

	robot := InitRobot(111, shared.Map{
		ExploredPath: make([]shared.PointStruct, 0),
		FrameOfRef:   1,
	})

	robotRPC := shared.RobotRPC{}
	rpc.Register(robotRPC)
	go rpc.Accept(listener)
	fmt.Println("Robot listening on port" + string(IPAddr))
	// for {
	// 	// wait for user input
	// 	// if button is pressed, break out of the loop
	// 	break
	// }

	for {
		// asynchronously check for other robots
		// if a robot is nearby, get IP address and make RPC call
		go robot.RespondToButtons()
		robot.Explore()
		break
	}

}

func InitRobot(rID uint, initMap shared.Map) *shared.RobotStruct {
	newRobot := shared.RobotStruct{
		RobotID:      rID,
		RMap:         initMap,
		JoiningSig:   make(chan bool),
		BusySig:      make(chan bool),
		WaitingSig:   make(chan bool),
		FreeSpaceSig: make(chan bool),
		WallSig:      make(chan bool),
		WalkSig:      make(chan bool),
	}
	// newRobot.CurPath.ListOfPCoordinates = append(newRobot.CurPath.ListOfPCoordinates, shared.PointStruct{PointKind: true})
	return &newRobot
}
