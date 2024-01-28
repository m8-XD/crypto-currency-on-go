package utils

import (
	"blockchain/pkg/entity"
	"bytes"
	"fmt"
)

const BUFFER_SIZE int = 10240

func fatal(err error, c *entity.Client) {
	fmt.Println(err)
	//TODO add ui pop-up
	c.Stop()
}

func Write(msg string, c *entity.Client) {
	peers := c.WritePeers()
	for _, peer := range peers {
		if peer == nil {
			continue
		}
		peer.Write([]byte(msg))
	}
}

func Read(c *entity.Client) []string {
	peers := c.ReadPeers()
	msgs := make([]string, len(peers))
	for i, peer := range peers {
		if peer == nil {
			continue
		}
		buff := make([]byte, BUFFER_SIZE)
		peer.Read(buff)
		msgs[i] = TrimAndCast(buff)
	}
	return msgs
}

func TrimAndCast(buff []byte) string {
	return string(bytes.Trim(buff, "\x00"))
}
