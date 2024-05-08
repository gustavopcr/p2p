package main

import (
	"io"
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
	// Handle incoming request
    // For simplicity, just send a byte slice back to the client
    for{
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil{
			if err == io.EOF{
				fmt.Println("EOF")
			}else{
				fmt.Println("error reading from conn: ", err)
			}
			if conn != nil{
				conn.Close()
			}
			return
		}
		fmt.Printf("data read from conn: %s:\n", string(buf))
	}
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