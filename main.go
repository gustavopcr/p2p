package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	//cli
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			continue
		}

		input = strings.TrimSpace(input)
		in := strings.Split(input, " ")
		if len(in) < 2 {
			fmt.Println("Error reading input: < 2")
			continue
		}
		command := in[0]
		switch command {
		case "send":
			filename := in[1]
			p.UploadFile(filename, sendChannel)

		}
	}
}
