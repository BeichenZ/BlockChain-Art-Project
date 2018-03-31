package main

import (
	"fmt"
	"net"
	"os"
)

func handleRequest(conn net.Conn) {
   buf := make([]byte, 1024)
   infoLen, err := conn.Read(buf)

   if err != nil {
	fmt.Println("ERROR")
   }

   fmt.Println(string(buf[:infoLen]))

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
		go handleRequest(conn)
	}
}
