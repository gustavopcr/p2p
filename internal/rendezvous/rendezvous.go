package rendezvous

import (
	"fmt"
	"net"

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
	tmpBuffer := make([]byte, 1024)
	n, _, err := p.ReadData(tmpBuffer)
	if err != nil {
		fmt.Println("Error receiving peer info:", err)
		return err
	}
	peerAddrStr := string(tmpBuffer[:n])
	peerAddr, err := net.ResolveUDPAddr("udp", peerAddrStr)
	if err != nil {
		fmt.Println("Error resolving peer address:", err)
		return err
	}

	p.PeersAddr = append(p.PeersAddr, peerAddr)
	return nil
}

func ConnectToPeers(server string) *peer.Peer { //substituir server
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

	return p
}
