package main

import (
	"fmt"
	"github.com/gustavopcr/p2p/internal/server"
)

func main() {
	//server.
	go server.StartServer()
	fmt.Println("p2p running...")
	
	for {
	}
	
}

