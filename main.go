package main

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"./shared"
	"github.com/DistributedClocks/GoVector/govec"
)

// TODO: Include golang GPIO

func main() {
	gob.Register(&net.TCPAddr{})
	gob.Register(&shared.TaskPayload{})

	// myIPAddr := GetLocalIP().String()
	// size, _ := net.ParseIP(GetLocalIP().String()).DefaultMask().Size()
	fmt.Println(GetLocalIP().String())
	fmt.Println("--------------------")
	ipv4Addr, ipv4Net, _ := net.ParseCIDR(GetLocalIP().String())
	fmt.Println(ipv4Addr)
	fmt.Println(ipv4Net)
	fmt.Println("----------------------")
	var ips []string
	for ip := ipv4Addr.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	fmt.Println(ips[1 : len(ips)-1])
	broad, _ := lastAddr(GetLocalIP())
	// pinger, err := ping.NewPinger(broad.String())
	// if err != nil {
	// 	panic(err)
	// }
	// pinger.Count = 3
	// pinger.Run()                 // blocks until finished
	// stats := pinger.Statistics() // get send/receive/rtt stats
	// fmt.Println(stats)
	fmt.Println(broad)
	// ones, _ := net.ParseIP(myIPAddr).DefaultMask().Size()
	// sub := ipsubnet.SubnetCalculator(myIPAddr, ones)
	// boardCastIP := sub.GetBroadcastAddress()
	// fmt.Println(boardCastIP)
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

func lastAddr(n *net.IPNet) (net.IP, error) { // works when the n is a prefix, otherwise...
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses.")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
