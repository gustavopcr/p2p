package main

import (
	"time"

	"github.com/gustavopcr/p2p/internal/rendezvous"
)

func main() {
	p := rendezvous.ConnectToPeers("localhost:8080")

	sendChannel := make(chan []byte, 10)
	receiveChannel := make(chan []byte, 10)
	//messageChannel := make(chan peer.Message)

	defer close(sendChannel)
	defer close(receiveChannel)

	//go peer.HandleMessage(messageChannel)
	go p.SendMessages(sendChannel)
	go p.DownloadFile(receiveChannel)
	p.UploadFile("alo.txt", sendChannel)
	time.Sleep(time.Second * 5)
	p.UploadFile("ola.txt", sendChannel)

	select {}
	/*
		subscribeToServer(ip)

		for{
			listen for events such as sending data another peer requested
		}
	*/
}
