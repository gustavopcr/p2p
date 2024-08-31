package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

type Message struct {
	MessageType    int
	SequenceNumber int
	Payload        []byte
}

func sendMessages(conn *net.UDPConn, peerAddr *net.UDPAddr, sendChannel <-chan []byte) error {
	for data := range sendChannel {
		_, err := conn.WriteToUDP(data, peerAddr)

		if err != nil {
			fmt.Println("Error sending message:", err)
			return err
		}
	}
	return nil
}

func (p *Peer) UploadFile(filename string) {
	sendChannel := make(chan []byte)
	defer close(sendChannel)
	for _, peerAddr := range p.PeersAddr {
		go sendMessages(p.conn, peerAddr, sendChannel)
	}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	tmpBuffer := make([]byte, 1024)
	for { // lendo arquivo
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)

		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := Message{MessageType: 1, SequenceNumber: 0, Payload: tmpBuffer[:n]}

				err = encoder.Encode(message)
				if err != nil {
					fmt.Println("Error encoding struct:", err)
					return
				}
				sendChannel <- buffer.Bytes()
				break
			}
			panic(err) // Handle other potential errors
		}
		message := Message{MessageType: 0, SequenceNumber: 1, Payload: tmpBuffer[:n]}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
	}
}

func (p *Peer) DownloadFile(messageChannel chan<- Message) {
	var msg Message
	tmpBuffer := make([]byte, 2048)
	for {
		var buffer bytes.Buffer
		decoder := gob.NewDecoder(&buffer)

		n, _, err := p.ReadData(tmpBuffer)
		fmt.Println("tmpBuffer: ", tmpBuffer)
		if err != nil {
			panic(err)
		}
		buffer.Write(tmpBuffer[:n])
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("err: ", err)
			panic(err)
		}
		messageChannel <- msg
	}
}

func HandleMessage(messageChannel <-chan Message) {
	for msg := range messageChannel {
		fmt.Println("msg: ", msg)
	}
}
