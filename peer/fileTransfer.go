package peer

import (
	"bufio"
	"io"
	"os"
)

func (p *Peer) SendFile(filename string) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(fi)

	for {
		n, err := r.Read(p.Buffer)
		if err != nil {
			if err == io.EOF {
				// No more data to read
				break
			}
			panic(err) // Handle other potential errors
		}

		// Process the bytes read
		//fmt.Print("arquivo enviado", string(p.Buffer[:n]))
		for _, peer := range p.PeersAddr {
			p.SendData(p.Buffer[:n], peer)
		}
	}
}
