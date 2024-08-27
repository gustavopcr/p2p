package rendezvous

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gustavopcr/p2p/peer"
)

func requestRendezvous(p *peer.Peer) error {
	message := []byte("Connect me")
	_, err := p.SendData(message, p.NexusAddr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	fmt.Println("Connection request sent to server")
	return nil
}

func addPeer(p *peer.Peer) error {
	//nexusConn.SetReadDeadline(time.Now().Add(10 * time.Second)) // Set a deadline for reading
	n, _, err := p.ReadData()
	if err != nil {
		fmt.Println("Error receiving peer info:", err)
		return err
	}
	peerAddrStr := string(p.Buffer[:n])
	peerAddr, err := net.ResolveUDPAddr("udp", peerAddrStr)
	if err != nil {
		fmt.Println("Error resolving peer address:", err)
		return err
	}

	p.PeersAddr = append(p.PeersAddr, peerAddr)
	return nil
}

func ConnectToPeers(server string) { //substituir server
	p, err := peer.NewPeer(server)
	if err != nil {
		panic(err)
	}

	err = requestRendezvous(p)
	if err != nil {
		panic(err)
	}
	// Receive peer information from the server
	err = addPeer(p)
	if err != nil {
		panic(err)
	}

	// Start sending packets to the peer using WriteToUDP
	go func() {
		for {
			fmt.Println("Sending data to peer...")
			for _, peerAddr := range p.PeersAddr {
				_, err := p.SendData([]byte(fmt.Sprintf("Hello Peer!: %d", rand.Int())), peerAddr)
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
		n, addr, err := p.ReadData()
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		fmt.Printf("Received message from %s: %s\n", addr, string(p.Buffer[:n]))
	}
}
