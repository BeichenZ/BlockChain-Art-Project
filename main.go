package main

import (
	"net/rpc"
	"./shared"
	"net"
	"log"
)

// TODO: Include golang GPIO

func main() {
	/// Need to change to different ip address. May to use a different library due to ad-hoc
	resolvedIPAddress, error := net.ResolveTCPAddr("tcp", "localhost")
	if error != nil {
		log.Fatal("Unable to resolve IP Address", error)
	}

	listener, error := net.ListenTCP("tcp", resolvedIPAddress)

	if error != nil {
		log.Fatal("Unable to create a listner", error)
	}

	robotRPC := shared.RobotRPC{}
	rpc.Register(robotRPC)
	go rpc.Accept(listener)

	for {
		// wait for user input
		// if button is pressed, break out of the loop
		break
	}

	for {
		// asynchronously check for other robots
		// if a robot is nearby, get IP address and make RPC call
		break

	}


}
