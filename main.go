package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"./GoVector/govec"
	"./shared"
)

// TODO: Include golang GPIO

func main() {
	gob.Register(&net.TCPAddr{})
	gob.Register(&shared.TaskPayload{})

	/// Need to change to different ip address. May to use a different library due to ad-hoc
	IPAddr := os.Args[1]
	RobotID, _ := strconv.Atoi(os.Args[2])
	Logger := govec.InitGoVector("Robot"+IPAddr, "LogFile"+IPAddr, true)
	resolvedIPAddr := IPAddr
	// resolvedIPAddress, error := net.ResolveTCPAddr("tcp", IPAddr)
	// if error != nil {
	// 	log.Fatal("Unable to resolve IP Address", error)
	// }

	robot := shared.InitRobot(RobotID, shared.Map{
		ExploredPath: make(map[shared.Coordinate]shared.PointStruct),
		FrameOfRef:   1,
	}, Logger)

	robotRPC := &shared.RobotRPC{PiRobot: robot}
	rpc.Register(robotRPC)
	listener, error := net.Listen("tcp", resolvedIPAddr)

	if error != nil {
		log.Fatal("Unable to create a listner", error)
	}
	go rpc.Accept(listener)
	fmt.Println("Robot listening on port " + string(IPAddr))
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
