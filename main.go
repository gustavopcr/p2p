package main

import (
	"github.com/gustavopcr/p2p/internal/rendezvous"
)

func main() {
	p := rendezvous.ConnectToPeers("localhost:8080")
	receiveChannel := make(chan []byte)

	go func() {
		p.UploadFile("alo.txt")
	}()

	for {
		p.DownloadFile(receiveChannel)
	}

	/*
		subscribeToServer(ip)

		for{
			listen for events such as sending data another peer requested
		}
	*/
}
