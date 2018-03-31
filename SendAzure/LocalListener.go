package main

import (
	"os/exec"

	"../shared"

	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type MapRPC struct {
	latestMergedMap shared.Map
}

func DecodeMap(sentMap []byte) shared.Map {

	//robotLogContent, _ := ioutil.ReadFile("./" + r.Logname)
	buf := bytes.NewBuffer(sentMap)

	var decodedMap shared.Map

	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	return decodedMap
}

func handleRequest(conn net.Conn) {

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Turn off Wi-Fi interface
	output, err := exec.Command("networksetup", "-setairportpower", "en0", "off").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	// Switch to UBC secure
	output, err = exec.Command("networksetup", "-setairportpower", "en0", "on").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

	//Wait until connected to Internet
WaitingForInternet:
	for {
		conn, err := net.Dial("tcp", "13.93.181.35:8080")
		if err != nil {
			continue
		} else {
			conn.Write(buf)
			break WaitingForInternet
		}
	}

	//Send Map to Azure

	//Change back to Ad-hoc
	output, err = exec.Command("networksetup", "-setairportnetwork", "en0", "PiAdHocNetwork").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))

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
