package peer

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/gustavopcr/p2p/constants"
	"github.com/gustavopcr/p2p/internal/file"
	"google.golang.org/protobuf/proto"
)

const (
	messageInit    = 0
	messagePayload = 1
	messageEOF     = 2
)

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

func (p *Peer) HandleMessages(messageChannel <-chan *Packet) {
	m := make(map[string]*file.FileManager)

	for msg := range messageChannel {
		switch msg.MessageType {
		case messageInit:
			fmt.Println("Init")
			f, err := file.NewFileManager(string(msg.Payload))
			if err != nil {
				panic(err)
			}
			m[msg.FileId] = f
		case messagePayload:
			fmt.Println("Payload")
			m[msg.FileId].File.WriteAt(msg.Payload, msg.Offset)
		case messageEOF:
			fmt.Println("EOF")
			if m[msg.FileId] != nil && m[msg.FileId].File != nil {
				err := m[msg.FileId].File.Close()
				if err != nil {
					panic(err)
				}
				delete(m, msg.FileId)
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

		message := &Packet{FileId: fileId.String(), MessageType: messageInit, SequenceNumber: 0, Offset: offset, Payload: []byte(fmt.Sprintf("p2p_%s", filename))}
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			panic(err)
		}
		sendChannel <- msgBytes
	}()

	for { // reading file
		n, err := r.Read(tmpBuffer)
		if err != nil {
			if err == io.EOF {
				message := &Packet{FileId: fileId.String(), MessageType: messageEOF, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
				msgBytes, err := proto.Marshal(message)
				if err != nil {
					panic(err)
				}
				sendChannel <- msgBytes
				break
			}
			panic(err)
		}
		message := &Packet{FileId: fileId.String(), MessageType: messagePayload, SequenceNumber: 0, Offset: offset, Payload: tmpBuffer[:n]}
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			panic(err)
		}
		sendChannel <- msgBytes
		offset += int64(n)
	}
}

func (p *Peer) DownloadFile(messageChannel chan<- *Packet) {
	for {

		tmpBuffer := make([]byte, constants.BufferSize)

		for {
			_, _, err := p.ReadData(tmpBuffer)
			if err != nil {
				if err == io.EOF {
					continue
				}
				panic(err)
			}
			message := &Packet{}
			if err := proto.Unmarshal(tmpBuffer, message); err != nil {
				panic(err)
			}
			messageChannel <- message
		}
	}
}
