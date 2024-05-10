package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	_ "time"
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
		fmt.Printf("data read from conn: %s:\n", string(buf))
		
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

func startServer() {
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
		go func(conn net.Conn){
			defer conn.Close()
			bufio := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			handleRequest(bufio)
		}(conn)
	}
}

func connectToServer(ip, port, filename string) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	conn.Write([]byte(filename))
	return conn, nil
}


func main() {
	go startServer()
	fmt.Println("p2p running...")
	for {
		/*
		conn, err := connectToServer("172.25.30.133", "6060", "HorasExtras.pdf")
		if err != nil {
			fmt.Println("error: ", err)
			conn.Close()
			continue
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
	*/
	}
}

