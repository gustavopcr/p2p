package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Create the address to send to
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error dialing UDP connection:", err)
		return
	}
	defer conn.Close()

	// Message to send
	message := []byte("Hello, UDP Server!")

	// Send the message to the server
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Println("Message sent to UDP server:", string(message))

	// Optional: wait for a response
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set a deadline for reading
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("No response from server or error reading response:", err)
	} else {
		fmt.Printf("Response from server: %s\n", string(buffer[:n]))
	}
}
