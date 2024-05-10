package main

import (
	"fmt"
	"time"
	"github.com/gustavopcr/p2p/internal/client"
)

func main() {
	for {
		conn, err := client.ConnectToServer("172.25.30.133", "6060", "HorasExtras.pdf")
		if err != nil {
			fmt.Println("error: ", err)
			conn.Close()
			return
		}
		for {
			tmp := make([]byte, 256)
			n, _ := conn.Read(tmp)
			if n > 0 {
				fmt.Println("tmp: ", string(tmp))
			} else {
				break
			}
		}
		conn.Close()
		time.Sleep(3 * time.Second)
	}
}

