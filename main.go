package main

import (
	"github.com/gustavopcr/p2p/internal/rendezvous"
	"github.com/gustavopcr/p2p/peer"
)

func main() {
	p := rendezvous.ConnectToPeers("localhost:8080")
	messageChannel := make(chan peer.Message)

	go p.UploadFile("alo.txt")

	//go peer.HandleMessage(messageChannel)

	for {
		p.DownloadFile(messageChannel)
	}

	/*
		subscribeToServer(ip)

		for{
			listen for events such as sending data another peer requested
		}
	*/
}
