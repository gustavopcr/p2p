package peer

import (
	"fmt"
	"net"
)

type Peer struct {
	Conn      *net.UDPConn
	NexusAddr *net.UDPAddr
	PeersAddr []*net.UDPAddr
}

func NewPeer(server string) (*Peer, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", server) // Public IP of the rendezvous server
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return nil, err
	}
	// Use a single connection for both server and Peer communication
	nexusConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return nil, err
	}
	return &Peer{Conn: nexusConn, NexusAddr: serverAddr, PeersAddr: make([]*net.UDPAddr, 0)}, nil
}

func (p *Peer) SendData(data []byte, PeerAddr *net.UDPAddr) (int, error) {
	n, err := p.Conn.WriteToUDP(data, PeerAddr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return 0, err
	}
	return n, nil
}

func (p *Peer) ReadData(buffer []byte) (int, *net.UDPAddr, error) {
	//nexusConn.SetReadDeadline(time.Now().Add(10 * time.Second)) // Set a deadline for reading
	n, addr, err := p.Conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error receiving Peer info:", err)
		return 0, nil, err
	}
	return n, addr, nil
}
