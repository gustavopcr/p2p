package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

type Message struct {
	MessageType    int
	SequenceNumber int
	Payload        []byte
}

func (p *Peer) SendMessages(sendChannel <-chan []byte) {
	for data := range sendChannel {
		for _, peerAddr := range p.PeersAddr {
			_, err := p.Conn.WriteToUDP(data, peerAddr)
			if err != nil {
				fmt.Println("Error sending message:", err)
				panic(err)
			}
		}
	}
}

func (p *Peer) UploadFile(filename string, sendChannel chan<- []byte) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	tmpBuffer := make([]byte, 1024)
	for { // lendo arquivo
		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := Message{MessageType: 1, SequenceNumber: 0, Payload: tmpBuffer[:n]}
				err = encoder.Encode(message)
				if err != nil {
					fmt.Println("Error encoding struct:", err)
					panic(err)
				}
				sendChannel <- buffer.Bytes()
				buffer.Reset()
				break
			}
			panic(err)
		}
		message := Message{MessageType: 0, SequenceNumber: 1, Payload: tmpBuffer[:n]}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
		buffer.Reset()
	}

}

func (p *Peer) DownloadFile(receiveChannel chan<- []byte) {
	var buffer bytes.Buffer
	tmpBuffer := make([]byte, 1024)
	for {
		n, _, err := p.ReadData(tmpBuffer)
		if err != nil {
			panic(err)
		}
		buffer.Write(tmpBuffer[:n])
		receiveChannel <- buffer.Bytes()
		buffer.Reset()
	}
}
