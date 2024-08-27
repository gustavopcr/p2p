package rendezvous

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func requestRendezvous(server string) (*net.UDPConn, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", server) // Public IP of the rendezvous server
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil, err
	}
	// Use a single connection for both server and peer communication
	nexusConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return nil, err
	}
	// Send a connection request to the server
	message := []byte("Connect me")
	_, err = nexusConn.WriteToUDP(message, serverAddr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return nil, err
	}
	fmt.Println("Connection request sent to server")
	return nexusConn, nil
}

func listenPeer(nexusConn *net.UDPConn, buffer []byte) (*net.UDPAddr, error) {
	//nexusConn.SetReadDeadline(time.Now().Add(10 * time.Second)) // Set a deadline for reading
	n, _, err := nexusConn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error receiving peer info:", err)
		return nil, err
	}

	peerAddrStr := string(buffer[:n])
	fmt.Println("Received peer info:", peerAddrStr)

	peerAddrStr = strings.TrimPrefix(peerAddrStr, "Peer: ")
	peerAddr, err := net.ResolveUDPAddr("udp", peerAddrStr)
	if err != nil {
		fmt.Println("Error resolving peer address:", err)
		return nil, err
	}

	fmt.Printf("Local address: %s\n", nexusConn.LocalAddr().String())
	return peerAddr, nil
}

func ConnectToPeers(server string) { //substituir server
	nexusConn, err := requestRendezvous(server)
	if err != nil {
		panic(err)
	}
	// Receive peer information from the server
	buffer := make([]byte, 1024)
	peerAddr, err := listenPeer(nexusConn, buffer)
	if err != nil {
		panic(err)
	}

	// Start sending packets to the peer using WriteToUDP
	go func() {
		for {
			fmt.Println("Sending data to peer...")
			_, err := nexusConn.WriteToUDP([]byte(fmt.Sprintf("Hello Peer!: %d", rand.Int())), peerAddr)
			if err != nil {
				fmt.Println("Error sending packet to peer:", err)
			} else {
				fmt.Println("Packet sent to peer.")
			}
			time.Sleep(3 * time.Second)
		}
	}()

	// Listen for incoming packets from the peer
	for {
		fmt.Println("Waiting to read data from peer...")
		n, addr, err := nexusConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		fmt.Printf("Received message from %s: %s\n", addr, string(buffer[:n]))
	}
}
