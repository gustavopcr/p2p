package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Message struct {
	MessageType    int
	SequenceNumber int
	Payload        []byte
}

func listenForMessages(conn *net.UDPConn, receiveChannel chan<- []byte) {
	buffer := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from UDP: ", err)
		}
		receiveChannel <- buffer
	}
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

func handleMessages(messageChannel <-chan []Message) {
	for msg := range messageChannel {
		fmt.Println("msg: ", msg)
	}
}

func (p *Peer) UploadFile(filename string) {

	sendChannel := make(chan []byte)
	for _, peerAddr := range p.PeersAddr {
		go sendMessages(p.conn, peerAddr, sendChannel)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	/*
		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}
		fileSize := fi.Size()
		_ = fileSize / int64(len(p.Buffer))
	*/
	r := bufio.NewReader(f)

	for { // lendo arquivo
		n, err := r.Read(buffer.Bytes())
		if err != nil {
			if err == io.EOF {
				message := Message{MessageType: 1, SequenceNumber: 0, Payload: p.Buffer[:n]}
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

		message := Message{MessageType: 0, SequenceNumber: 1, Payload: p.Buffer[:n]}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
	}
}

func (p *Peer) DownloadFile() {
	for n, addr, err := p.ReadData(); n >= 0; n, addr, err = p.ReadData() {
		if err != nil {
			panic(err)
		}
		fmt.Println("reading: ", p.Buffer, " from: ", addr)
	}
}
