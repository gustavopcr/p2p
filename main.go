package main

import (
	"fmt"

	"github.com/gustavopcr/p2p/internal/rendezvous"
	"github.com/gustavopcr/p2p/peer"
)

func main() {
	p := rendezvous.ConnectToPeers("localhost:8080")

	sendChannel := make(chan []byte, 10)
	//receiveChannel := make(chan []byte, 10)
	messageChannel := make(chan peer.Message, 10)

	defer close(sendChannel)
	defer close(messageChannel)

	//go peer.HandleMessage(messageChannel)
	go p.SendMessages(sendChannel)
	go p.DownloadFile(messageChannel)
	p.UploadFile("alo.txt", sendChannel)
	fmt.Println("time.sleep...")
	p.UploadFile("ola.txt", sendChannel)
	for msg := range messageChannel {
		fmt.Println("msg.Payload: ", string(msg.Payload))
	}
}
