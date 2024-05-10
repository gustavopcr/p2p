package client

import (
	"fmt"
	"net"
)

func ConnectToServer(ip, port, filename string) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	conn.Write([]byte(filename))
	return conn, nil
}
