package main

import "github.com/gustavopcr/p2p/internal/rendezvous"

func main() {
	rendezvous.ConnectToPeers("localhost:8080")
	/*
		subscribeToServer(ip)

		for{
			listen for events such as sending data another peer requested
		}
	*/
}
