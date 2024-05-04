package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	data, err := os.ReadFile("teste3.7z")
	if err != nil {
		fmt.Printf("Erro ao ler arquivo: %s\n", err)
	}

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Printf("Error dialing TCP: %s\n", err)
		return
	}
	defer conn.Close()

	sendData(len(data)/2, &data, &conn)
}
