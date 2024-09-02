package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

type Message struct {
	FileID         uuid.UUID
	MessageType    int
	SequenceNumber int
	Offset         int
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

func (p *Peer) HandleMessages(messageChannel <-chan Message) {
	for msg := range messageChannel {
		fmt.Println("msg.Payload: ", string(msg.Payload))
	}
}

func (p *Peer) UploadFile(filename string, sendChannel chan<- []byte) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	tmpBuffer := make([]byte, 1024)

	fileId := uuid.New()
	offset := 0
	for { // lendo arquivo
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				/*
					type Message struct {
						MessageType    int
						SequenceNumber int
						Offset         uint64
						FileID         uuid.UUID
						Payload        []byte
					}
				*/
				message := Message{FileID: fileId, MessageType: 1, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
				err = encoder.Encode(message)
				if err != nil {
					fmt.Println("Error encoding struct:", err)
					panic(err)
				}
				sendChannel <- buffer.Bytes()
				break
			}
			panic(err)
		}
		message := Message{FileID: fileId, MessageType: 1, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
	}

}

func (p *Peer) DownloadFile(messageChannel chan<- Message) {
	tmpBuffer := make([]byte, 1024)
	for {
		var buffer bytes.Buffer
		decoder := gob.NewDecoder(&buffer)
		n, _, err := p.ReadData(tmpBuffer)
		if err != nil {
			panic(err)
		}
		buffer.Write(tmpBuffer[:n])
		var msg Message
		err = decoder.Decode(&msg)
		if err != nil {
			panic(err)
		}
		messageChannel <- msg
		buffer.Reset()
	}
}
