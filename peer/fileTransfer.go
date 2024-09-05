package peer

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/gustavopcr/p2p/constants"
	"github.com/gustavopcr/p2p/internal/file"
)

const (
	messageInit    = 0
	messagePayload = 1
	messageEOF     = 2
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
			fmt.Println("message sent to peerAddr: ", peerAddr)
			_, err := p.Conn.WriteToUDP(data, peerAddr)
			if err != nil {
				fmt.Println("Error sending message:", err)
				panic(err)
			}
		}
	}
}

func (p *Peer) HandleMessages(messageChannel <-chan Message) {
	m := make(map[uuid.UUID]*file.FileManager)

	for msg := range messageChannel {
		switch msg.MessageType {
		case messageInit:
			fmt.Println("Init")
			f, err := file.NewFileManager(string(msg.Payload))
			if err != nil {
				panic(err)
			}
			m[msg.FileID] = f
		case messagePayload:
			fmt.Println("Payload")
			m[msg.FileID].File.WriteAt(msg.Payload, msg.Offset)
		case messageEOF:
			fmt.Println("EOF")
			if m[msg.FileID] != nil && m[msg.FileID].File != nil {
				err := m[msg.FileID].File.Close()
				if err != nil {
					panic(err)
				}
				delete(m, msg.FileID)
			}
		}
	}
}

func (p *Peer) UploadFile(filename string, sendChannel chan<- []byte) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	tmpBuffer := make([]byte, 1000) // set size as 1000 to avoid encode/decode bug, malformed struct due to ignoring struct header size

	fileId := uuid.New()
	offset := int64(0)

	func() {
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		message := Message{FileID: fileId, MessageType: messageInit, SequenceNumber: 0, Offset: offset, Payload: []byte(fmt.Sprintf("p2p_%s", filename))}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			panic(err)
		}
		sendChannel <- buffer.Bytes()
	}()

	for { // lendo arquivo
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := Message{FileID: fileId, MessageType: messageEOF, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
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
		message := Message{FileID: fileId, MessageType: messagePayload, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding struct:", err)
			return
		}
		sendChannel <- buffer.Bytes()
		buffer.Reset()
		offset += int64(n)
	}
}

func (p *Peer) DownloadFile(messageChannel chan<- Message) {

	tmpBuffer := make([]byte, constants.BufferSize)

	for {
		var buffer bytes.Buffer
		decoder := gob.NewDecoder(&buffer)
		n, _, err := p.ReadData(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				continue
			}
			panic(err)
		}

		_, err = buffer.Write(tmpBuffer[:n])
		if err != nil {
			fmt.Println("Error writing to buffer:", err)
			panic(err)
		}
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
