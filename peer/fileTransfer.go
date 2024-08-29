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

func (p *Peer) UploadFile(filename string) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	fileSize := fi.Size()
	_ = fileSize / int64(len(p.Buffer))

	r := bufio.NewReader(f)

	for {
		n, err := r.Read(p.Buffer)
		if err != nil {
			if err == io.EOF {
				for _, peer := range p.PeersAddr {
					message := Message{MessageType: 1, SequenceNumber: 1, Payload: p.Buffer[:n]}
					err := encoder.Encode(message)
					if err != nil {
						fmt.Println("Error encoding struct:", err)
						return
					}
					p.SendData(buffer.Bytes(), peer)
				}
				break
			}
			panic(err) // Handle other potential errors
		}

		for _, peer := range p.PeersAddr {
			message := Message{MessageType: 0, SequenceNumber: 1, Payload: p.Buffer[:n]}
			err := encoder.Encode(message)
			if err != nil {
				fmt.Println("Error encoding struct:", err)
				return
			}
			p.SendData(buffer.Bytes(), peer)
		}
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
