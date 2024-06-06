package main

import (
	"fmt"
	"net/http"
	"github.com/gustavopcr/p2p/internal/server"
	"github.com/gustavopcr/p2p/types"
	"time"
	"encoding/json"
)

func main() {
	var peers = make([]types.Peer, 0)
	go server.StartServer()	
	fmt.Println("p2p running...")
	
	for {
		r, err := http.Get("http://localhost:8080/peers")
		if err != nil{
			panic(err)
		}
		err = json.NewDecoder(r.Body).Decode(&peers)
		fmt.Println(peers)
		time.Sleep(30 * time.Second)
	}
	
}

