package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)
type CoordinateStruct struct {
	X float64
	Y float64
	IsItFreeToRoam bool
}

type InformationToWebApp struct {
	sync.RWMutex
	all map[int][]CoordinateStruct
}

var allInfo = InformationToWebApp{all: make( map[int][]CoordinateStruct)}


func handleRequestFromLocalListener(conn net.Conn) {
   buf := make([]byte, 1024)
   infoLen, err := conn.Read(buf)

   if err != nil {
	fmt.Println("ERROR")
   }


	newMap := DecodeMap(buf[:infoLen])
	fmt.Println("Successfully decoded info")
	fmt.Println(newMap.ExploredPath)
	fmt.Println(newMap.FrameOfRef)

	tmpListOfCoordinates := make([]CoordinateStruct, len(newMap.ExploredPath))

	for cord, points := range  newMap.ExploredPath{
		tmpCordStruct := CoordinateStruct{}

		tmpCordStruct.X = cord.X
		tmpCordStruct.Y = cord.Y

		tmpCordStruct.IsItFreeToRoam = points.PointKind

		tmpListOfCoordinates = append(tmpListOfCoordinates, tmpCordStruct)
	}

	allInfo.Lock()
	allInfo.all[newMap.FrameOfRef] = tmpListOfCoordinates
	allInfo.Unlock()


}	

func main() {
	port := os.Args[1]
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on port", port)
	for {
		// Listen for an incoming connection.
		conn , err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		} 
		go handleRequestFromLocalListener(conn)
	}
}
