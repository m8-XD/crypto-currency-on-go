package utils

import (
	"blockchain/pkg/cryptography"
	"blockchain/pkg/entity"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const BUFFER_SIZE int = 10240

func Write(msg string, c *entity.Client) {
	peers := c.WritePeers()
	fmt.Println("write: " + msg)
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

func IsNumber(num string) bool {
	_, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return false
	}
	return true
}

func SendTX(c *entity.Client, kPair *cryptography.KeyPair, payload string) {
	privKey, err := kPair.Private()
	if err != nil {
		return
	}
	ds, err := cryptography.Sign(privKey, []byte(payload))
	if err != nil {
		fmt.Println("couldn't generate digital signature: " + err.Error())
		return
	}
	txText := strings.Join([]string{payload, ds}, ":")
	Write(txText, c)
}

func ChooseBlock(amount float64) string {
	return "someblock: " + fmt.Sprint(amount)
}