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

	//bgpio "./gpio"
	"./shared"
	"github.com/DistributedClocks/GoVector/govec"
)

const (
	FrontObstacleButton_Pin uint = 5
	FrontEmptyButton_Pin    uint = 6
	LeftObstacleButton_Pin  uint = 13
	RightObstacleButton_Pin uint = 19
)

// TODO: Include golang GPIO

func main() {
	gob.Register(&net.TCPAddr{})
	gob.Register(&shared.TaskPayload{})
	gob.Register(&shared.Neighbour{})

	/// Need to change to different ip address. May to use a different library due to ad-hoc
	Port := os.Args[1]
	RobotID, _ := strconv.Atoi(os.Args[2])
	RobotInitialPositionX, error := strconv.ParseFloat(os.Args[3], 64)
	RobotInitialPositionY, error := strconv.ParseFloat(os.Args[4], 64)

	if error != nil {
		fmt.Println("fail to parse the inital current location")
		os.Exit(1)
	}
	RobotInitialPosition := shared.Coordinate{RobotInitialPositionX, RobotInitialPositionY}
	// Logger := govec.InitGoVector("Port", "LogFile"+Port)

	fmt.Println("Robot IP Address:", GetLocalIP().String())
	ipv4Addr, ipv4Net, _ := net.ParseCIDR(GetLocalIP().String())
	fmt.Println("--------------------")
	fmt.Println(ipv4Addr)
	fmt.Println(ipv4Net)
	fmt.Println(ipv4Addr.String() + Port)
	fmt.Println("----------------------")

	Logger := govec.InitGoVector("Robot"+ipv4Addr.String()+Port, "LogFile"+ipv4Addr.String()+Port)
	resolvedIPAddr := Port
	// resolvedIPAddress, error := net.ResolveTCPAddr("tcp", Port)
	// if error != nil {
	// 	log.Fatal("Unable to resolve IP Address", error)
	// }
	logname := "Robot" + ipv4Addr.String() + Port + "-Log.txt"
	robot := shared.InitRobot(RobotID, shared.Map{
		ExploredPath: make(map[shared.Coordinate]shared.PointStruct),
		FrameOfRef:   1,
	}, RobotInitialPosition, Logger, ipv4Addr.String()+Port, logname)
	fmt.Println("Robot current location before reading from log ", robot.CurLocation)
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
	fmt.Println("Robot listening on port " + string(Port))

	// fmt.Println(logname)
	if _, err := os.Stat("./" + logname); os.IsNotExist(err) {
		logFile, e := robot.CreateLog()
		if e != nil {
			fmt.Println(e)
		}
		robotLog := robot.ProduceLogInfo()
		logInfo := robot.EncodeRobotLogInfo(robotLog)
		logFile.WriteString(logInfo)
		logFile.Close()
	} else {
		file, err := os.Stat("./" + robot.Logname)
		if err != nil {
			fmt.Println("error opening the file")
		}
		size := file.Size()
		if size != 0 {
			fmt.Println("info found in the log, reading from the log to revive robot...")
			robot.ReadFromLog()
		}
	}

	fmt.Println("AFTER reading from log: Current location", robot.CurLocation,
		"Current energy ", robot.RobotEnergy, "Current task ")
	var ips []string
	for ip := ipv4Addr.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// fmt.Println(ips[1 : len(ips)-1])
	ips = ips[1 : len(ips)-2]

	timeout := time.Duration(100 * time.Millisecond)
	go scanForNeighbours(ips[:5], ipv4Addr, timeout, robot, Port)
	go robot.CallNeighbours()
	//leftObstacleButtonPin := bgpio.NewInput(LeftObstacleButton_Pin)
	//rightObstacleButtonPin := bgpio.NewInput(RightObstacleButton_Pin)
	//frontEmptyButtonPin := bgpio.NewInput(FrontEmptyButton_Pin)
	//frontObstacleButtonPin := bgpio.NewInput(FrontObstacleButton_Pin)
	//
	//go robot.MonitorButtonsOnPins(leftObstacleButtonPin)
	//go robot.MonitorRightButtonOnPin(rightObstacleButtonPin)
	//go robot.MonitorFrontButtonOnPin(frontEmptyButtonPin)
	//go robot.MonitorFrontObsButtonOnPin(frontObstacleButtonPin)

	// robot.MonitorButtons()
	//go robot.SendMapToLocalServer()
	// for {
	// 	// wait for user input
	// 	// if button is pressed, break out of the loop
	// 	break
	// }

	// asynchronously check for other robots
	// if a robot is nearby, get IP address and make RPC call
	go robot.RespondToButtons()
	robot.Explore()

}

func scanForNeighbours(ips []string, ipv4Addr net.IP, timeout time.Duration, robot *shared.RobotStruct, Port string) {
	for {
		//fmt.Println("Looking for neighbours...")
		for _, ip := range ips {
			if ip == ipv4Addr.String() {
				continue
			}
			_, err := net.DialTimeout("tcp", ip+":5000", timeout)
			if err == nil {
				//log.Println("Able to locate neighbour")
				// Start registeration protocol
				// TODO -- maybe potential bug - hardcoded stuff
				robot.PossibleNeighbours.Add(ip + ":8080")
				//robot.PossibleNeighbours.Add(ip + Port)
				// robot.PossibleNeighbours = append(robot.PossibleNeighbours, ip+":5000")
				//fmt.Println(robot.PossibleNeighbours)
				neighbourIPAddr := ""
				client, err := rpc.Dial("tcp", ip+":5000")
				if err != nil {
					continue
					fmt.Println(err)
				}
				error := client.Call("RobotRPC.RegisterNeighbour", ipv4Addr.String()+Port, &neighbourIPAddr)
				if error != nil {
					fmt.Println(error.Error())
				} else {
					robot.PossibleNeighbours.Add(neighbourIPAddr)
					client.Close()
				}
			}
		}
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
