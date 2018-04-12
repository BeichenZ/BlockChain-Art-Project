package main

import (
	"os/exec"
	"fmt"
	"net"
	"os"
	"encoding/gob"
	"../shared"
	"bytes"
)

func handleRequest(conn net.Conn) {

	buf := make([]byte, 65536)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	buffer := bytes.NewBuffer(buf)

	var decodedMap shared.Map

	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(decodedMap)

	// Turn off Wi-Fi interface
	_, err = exec.Command("networksetup", "-setairportpower", "en0", "off").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println("Turned off the WIFI")

	// Switch to UBC secure
	_, err = exec.Command("networksetup", "-setairportpower", "en0", "on").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println("Turned on the WIFI")


	//Wait until connected to Internet
WaitingForInternet:
	for {
		//13.91.38.239:8080 is the Azure IP
		conn, err := net.Dial("tcp", "13.91.38.239:8080")
		if err != nil {
			continue
		} else {
			_,error := conn.Write(buf)
			if error != nil{
				fmt.Println("Error on sending the map to Azure")
			}
			break WaitingForInternet
		}
	}

	//Send Map to Azure
	fmt.Println("Successfully sent map to azure")

	//Change back to Ad-hoc
	_, err = exec.Command("networksetup", "-setairportnetwork", "en0", "PiAdHocNetwork").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println("Turned back to adhoc")

	// os.Exit(0)

}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on 169.254.82.142:8081")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			// os.Exit(1)
		}
		// Handle connections in a new goroutine.
		handleRequest(conn)
	}
}

/*
networksetup -setairportnetwork en0 PiAdHocNetwork
networksetup -setairportpower en0 off
networksetup -setairportpower en0 on
*/
