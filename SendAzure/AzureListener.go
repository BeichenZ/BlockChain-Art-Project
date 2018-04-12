package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"bytes"
	"encoding/gob"
	"../shared"
)

type CoordinateStruct struct {
	X              float64
	Y              float64
	IsItFreeToRoam bool
}

type InformationToWebApp struct {
	sync.RWMutex
	all map[int][]CoordinateStruct
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

var allInfo = InformationToWebApp{all: make(map[int][]CoordinateStruct)}

func GetAllMaps(w http.ResponseWriter, r *http.Request) {
	s, err := json.Marshal(allInfo.all)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(s)
}

func handleRequestFromLocalListener(conn net.Conn) {
	buf := make([]byte, 65536)
	infoLen, err := conn.Read(buf)

	if err != nil {
		fmt.Println("ERROR")
	}

	newMap := DecodeMap(buf[:infoLen])
	fmt.Println("Successfully decoded info")
	fmt.Println(newMap.ExploredPath)
	fmt.Println(newMap.FrameOfRef)

	tmpListOfCoordinates := make([]CoordinateStruct, len(newMap.ExploredPath))

	for cord, points := range newMap.ExploredPath {
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
	mux := http.NewServeMux()
	mux.HandleFunc("/getallmaps", GetAllMaps)
	go http.ListenAndServe(":5000", mux)

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
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequestFromLocalListener(conn)
	}
}
