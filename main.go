package main

import (
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
    defer conn.Close()
	// Handle incoming request
    // For simplicity, just send a byte slice back to the client
    data := []byte("Hello, peer!")
    conn.Write(data)
}


func startServer() {
	l, err := net.Listen("tcp", ":6060")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("go handleRequest(conn)")
		go handleRequest(conn)
	}
}

func main() {
	go startServer()
	fmt.Println("p2p running...")
	for{
		
	}
}