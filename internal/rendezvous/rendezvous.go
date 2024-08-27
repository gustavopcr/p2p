package rendezvous

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type peer struct {
	conn      *net.UDPConn
	nexusAddr *net.UDPAddr
	peersAddr []*net.UDPAddr
	buffer    []byte
}

func newPeer(server string) (*peer, error) {
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
	return &peer{conn: nexusConn, nexusAddr: serverAddr, peersAddr: make([]*net.UDPAddr, 0), buffer: make([]byte, 1024)}, nil
}

func (p *peer) requestRendezvous() error {
	message := []byte("Connect me")
	_, err := p.conn.WriteToUDP(message, p.nexusAddr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	fmt.Println("Connection request sent to server")
	return nil
}

func (p *peer) listenPeer() error {
	//nexusConn.SetReadDeadline(time.Now().Add(10 * time.Second)) // Set a deadline for reading
	n, _, err := p.conn.ReadFromUDP(p.buffer)
	if err != nil {
		fmt.Println("Error receiving peer info:", err)
		return err
	}

	peerAddrStr := string(p.buffer[:n])
	fmt.Println("Received peer info:", peerAddrStr)

	peerAddrStr = strings.TrimPrefix(peerAddrStr, "Peer: ")
	peerAddr, err := net.ResolveUDPAddr("udp", peerAddrStr)
	if err != nil {
		fmt.Println("Error resolving peer address:", err)
		return err
	}

	p.peersAddr = append(p.peersAddr, peerAddr)
	return nil
}

func ConnectToPeers(server string) { //substituir server
	p, err := newPeer(server)
	if err != nil {
		panic(err)
	}
	//err := requestRendezvous(server)
	err = p.requestRendezvous()
	if err != nil {
		panic(err)
	}
	// Receive peer information from the server
	err = p.listenPeer()
	if err != nil {
		panic(err)
	}

	// Start sending packets to the peer using WriteToUDP
	go func() {
		for {
			fmt.Println("Sending data to peer...")
			for _, peerAddr := range p.peersAddr {
				_, err := p.conn.WriteToUDP([]byte(fmt.Sprintf("Hello Peer!: %d", rand.Int())), peerAddr)
				if err != nil {
					fmt.Println("Error sending packet to peer:", err)
				} else {
					fmt.Println("Packet sent to peer.")
				}
				time.Sleep(3 * time.Second)
			}
		}
	}()

	// Listen for incoming packets from the peer
	for {
		fmt.Println("Waiting to read data from peer...")
		n, addr, err := p.conn.ReadFromUDP(p.buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		fmt.Printf("Received message from %s: %s\n", addr, string(p.buffer[:n]))
	}
}
