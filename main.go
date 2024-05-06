package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
    "syscall"
)

func handleRequest(conn net.Conn) {
    defer conn.Close()
	// Handle incoming request
    // For simplicity, just send a byte slice back to the client
    //data := []byte("Hello, peer!")
    //conn.Write(data)
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

	// Start CLI
	fmt.Println("CLI is running...")

	// Wait for interrupt signal (Ctrl+C) to exit the program
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Exiting...")
}