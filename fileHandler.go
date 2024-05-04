package main

import (
	"fmt"
	"net"
)

func HelloWorld() {
	fmt.Println("hello World")
}

func sendData(chunkSize int, arr *[]byte, con *(net.Conn)) {
	n := len(*arr)
	for i := 0; i < n; i += chunkSize {
		upper := i + chunkSize
		if upper > n {
			upper = n
		}
		go (*con).Write((*arr)[i:upper])
	}
}
