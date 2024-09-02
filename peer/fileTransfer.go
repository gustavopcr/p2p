package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/gustavopcr/p2p/internal/file"
)

type Message struct {
	FileID         uuid.UUID
	MessageType    int
	SequenceNumber int
	Offset         int64
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
	f, err := file.NewFileManager("testando.txt")
	if err != nil {
		panic(err)
	}
	for msg := range messageChannel {
		f.File.WriteAt(msg.Payload, msg.Offset)
		//fmt.Println("msg.Payload: ", string(msg.Payload))
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
	offset := int64(0)
	for { // lendo arquivo
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
<<<<<<< HEAD
		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := Message{FileID: fileId, MessageType: 1, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
=======

		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := Message{MessageType: 1, SequenceNumber: 0, Payload: tmpBuffer[:n]}

>>>>>>> 84b4dd4ac639d75e116535b1137cc424f0baa688
				err = encoder.Encode(message)
				if err != nil {
					fmt.Println("Error encoding struct:", err)
					panic(err)
				}
				sendChannel <- buffer.Bytes()
				offset += int64(n + 1)
				break
			}
			panic(err)
		}
<<<<<<< HEAD
		message := Message{FileID: fileId, MessageType: 1, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
=======
		message := Message{MessageType: 0, SequenceNumber: 1, Payload: tmpBuffer[:n]}
>>>>>>> 84b4dd4ac639d75e116535b1137cc424f0baa688
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
		offset += int64(n + 1)
	}
}

func (p *Peer) DownloadFile(messageChannel chan<- Message) {
<<<<<<< HEAD
	tmpBuffer := make([]byte, 1024)
	for {
		var buffer bytes.Buffer
		decoder := gob.NewDecoder(&buffer)
=======
	var msg Message
	tmpBuffer := make([]byte, 2048)
	for {
		var buffer bytes.Buffer
		decoder := gob.NewDecoder(&buffer)

>>>>>>> 84b4dd4ac639d75e116535b1137cc424f0baa688
		n, _, err := p.ReadData(tmpBuffer)
		if err != nil {
			panic(err)
		}
		buffer.Write(tmpBuffer[:n])
		var msg Message
		err = decoder.Decode(&msg)
		if err != nil {
			fmt.Println("err: ", err)
			panic(err)
		}
		messageChannel <- msg
		buffer.Reset()
	}
}
