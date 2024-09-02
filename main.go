package main

import (
	"github.com/gustavopcr/p2p/internal/rendezvous"
	"github.com/gustavopcr/p2p/peer"
)

func main() {
	p := rendezvous.ConnectToPeers("localhost:8080")

	sendChannel := make(chan []byte, 10)
	messageChannel := make(chan peer.Message, 10)

	defer close(sendChannel)
	defer close(messageChannel)

	go p.SendMessages(sendChannel)
	go p.DownloadFile(messageChannel)
	go p.HandleMessages(messageChannel)

	p.UploadFile("alo.txt", sendChannel)

	select {}
}
