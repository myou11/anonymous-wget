package main

import (
	"fmt"
	"net"
	"os"
	"proj3/ssnet"
	"strconv"
)

func checkErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", e.Error())
		os.Exit(1)
	}
}

var (
	port           uint16
	listener       net.Listener
	portDefaultMsg string
)

func init() {
	portDefaultMsg = "port to bind the server at. Defaults to 8989"
	port = 0
}

func cleanup() {

}

func main() {
	//port = 8991
	//fmt.Println(port)

	// use user-specified port
	if len(os.Args) > 1 {
		port64, err := strconv.ParseUint(os.Args[1], 10, 16)
		checkErr(err)
		port = uint16(port64)
	}

	hostname, err := os.Hostname()
	checkErr(err)
	localhostPort := fmt.Sprintf("%s:%d", hostname, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", localhostPort)
	checkErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	defer listener.Close()
	listenerAddr := listener.Addr()
	fmt.Printf("listening on %s\n", listenerAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Don't want to exit program if a connection couldn't be accepted
			// Keep listening for connections
			continue
		}
		fmt.Println("ss: conn accepted")
		go ssnet.HandleReqFromConn(conn)
	}
}
