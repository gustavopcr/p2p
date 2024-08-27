package peer

import (
	"fmt"
	"net"
)

type Peer struct {
	conn      *net.UDPConn
	NexusAddr *net.UDPAddr
	PeersAddr []*net.UDPAddr
	Buffer    []byte
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
	return &Peer{conn: nexusConn, NexusAddr: serverAddr, PeersAddr: make([]*net.UDPAddr, 0), Buffer: make([]byte, 1024)}, nil
}

func (p *Peer) SendData(data []byte, PeerAddr *net.UDPAddr) (int, error) {
	n, err := p.conn.WriteToUDP(data, PeerAddr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return 0, err
	}
	fmt.Println("Connection request sent to server")
	return n, nil
}

func (p *Peer) ReadData() (int, *net.UDPAddr, error) {
	//nexusConn.SetReadDeadline(time.Now().Add(10 * time.Second)) // Set a deadline for reading
	n, addr, err := p.conn.ReadFromUDP(p.Buffer)
	if err != nil {
		fmt.Println("Error receiving Peer info:", err)
		return 0, nil, err
	}
	return n, addr, nil
}
