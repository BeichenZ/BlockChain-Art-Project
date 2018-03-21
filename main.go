package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"./shared"
	"github.com/DistributedClocks/GoVector/govec"
)

// TODO: Include golang GPIO

func main() {
	gob.Register(&net.TCPAddr{})
	gob.Register(&shared.TaskPayload{})

	/// Need to change to different ip address. May to use a different library due to ad-hoc
	IPAddr := os.Args[1]
	RobotID, _ := strconv.Atoi(os.Args[2])
	Logger := govec.InitGoVector("IPAddr", "LogFile"+IPAddr)
	resolvedIPAddr := IPAddr
	// resolvedIPAddress, error := net.ResolveTCPAddr("tcp", IPAddr)
	// if error != nil {
	// 	log.Fatal("Unable to resolve IP Address", error)
	// }

	robot := shared.InitRobot(RobotID, shared.Map{
		ExploredPath: make([]shared.PointStruct, 0),
		FrameOfRef:   1,
	}, Logger)

	// Open up user defined port RPC connection
	robotRPC := &shared.RobotRPC{PiRobot: robot}
	rpc.Register(robotRPC)
	listener, error := net.Listen("tcp", resolvedIPAddr)
	if error != nil {
		log.Fatal("Unable to create a listner", error)
	}

	// Open up port 5000 for broadcasting
	registerListener, error := net.Listen("tcp", ":5000")
	if error != nil {
		log.Fatal("Unable to create a listner", error)
	}
	go rpc.Accept(listener)
	go rpc.Accept(registerListener)
	fmt.Println("Robot listening on port " + string(IPAddr))

	fmt.Println("Robot IP Address:", GetLocalIP().String())
	ipv4Addr, ipv4Net, _ := net.ParseCIDR(GetLocalIP().String())
	fmt.Println("--------------------")
	fmt.Println(ipv4Addr)
	fmt.Println(ipv4Net)
	fmt.Println("----------------------")
	var ips []string
	for ip := ipv4Addr.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	fmt.Println(ips[1 : len(ips)-1])
	ips = ips[1 : len(ips)-1]

	timeout := time.Duration(10 * time.Millisecond)
	for _, ip := range ips {
		_, err := net.DialTimeout("tcp", ip+":5000", timeout)
		if err != nil {
			log.Println("Site unreachable, error: ", err)
		}
	}

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

func GetLocalIP() *net.IPNet {
	addrs, _ := net.InterfaceAddrs()

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet
			}
		}
	}
	return nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
