package rendezvous

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func handleRequest(rw *bufio.ReadWriter) {
	for {
		buf := make([]byte, 256)
		n, err := rw.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
			} else {
				fmt.Println("error reading from conn: ", err)
			}
			return
		}
		fmt.Printf("data read from conn: %s\n", string(buf))

		file, err := os.Open(string(buf[:n]))
		if err != nil {
			fmt.Println("err file: ", err)
			return
		}
		defer file.Close()
		bufio := bufio.NewReader(file)

		for {
			tmp := make([]byte, 256)
			n, err := bufio.Read(tmp)
			if err != nil {
				if err == io.EOF {
					fmt.Println("EOF")
				} else {
					fmt.Println("error reading from conn: ", err)
				}
				return
			}
			if n > 0 {
				_, err := rw.Write(tmp[:n])
				if err != nil {
					fmt.Println("error writing to conn: ", err)
					return
				}
				rw.Flush()
			}
		}
	}
}

func StartServer() {
	l, err := net.Listen("tcp", ":6060")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error")
		}
		go func(conn net.Conn) {
			defer conn.Close()
			bufio := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			handleRequest(bufio)
		}(conn)
	}
}
