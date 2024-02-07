package mining

import (
	"blockchain/pkg/entity"
	"blockchain/pkg/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Miner struct {
	client   *entity.Client
	chain    []node
	chainMut sync.Mutex
}

// transaction
type tx struct {
	WAddr     string
	RecWAddr  string //reciever wallet address
	Amount    float64
	Change    float64
	BHash     string //block hash
	Timestamp int64
	DS        string //digital signature
	Payload   string
}

func (m *Miner) IsRunning() bool {
	return m.client.IsRunning()
}

func (m *Miner) Client() *entity.Client {
	return m.client
}

func (m *Miner) SetClient(c *entity.Client) {
	m.client = c
}

func (m *Miner) Start() {
	m.client.Start()
	m.chain = make([]node, 0)
	m.chain = append(m.chain, node{
		Header: "0x00000000000000000000000000000000000000001",
		TX: tx{
			RecWAddr: "0x4e2448db9e9b30178433434ccd8ea9534b5e28e2",
			Amount:   123,
		},
	})
}

func (m *Miner) AddTX(txRaw string) {
	tx, err := parseTX(txRaw)
	if err == nil {
		go Mine(m, tx)
	} else {
		fmt.Println(err)
	}
}

func (m *Miner) CopyChain() []node {
	chainCopy := make([]node, len(m.chain))
	copy(chainCopy, m.chain)
	return chainCopy
}

func parseTX(txRaw string) (txn tx, err error) {
	parsedTX := strings.Split(txRaw, ":")
	if len(parsedTX) != 2 {
		err = errors.New("invalid TX")
		return
	}
	payloadRaw, ds := parsedTX[0], parsedTX[1]
	dsBytes, err := base64.StdEncoding.DecodeString(ds)
	ds = string(dsBytes)
	if err != nil {
		return
	}
	payload := strings.Split(payloadRaw, ",")
	amount, err := strconv.ParseFloat(payload[2], 64)
	if err != nil {
		return
	}
	change, err := strconv.ParseFloat(payload[3], 64)
	if err != nil {
		return
	}
	timestamp, err := strconv.ParseInt(payload[5], 10, 64)
	if err != nil {
		return
	}
	txn = tx{
		WAddr:     payload[0],
		RecWAddr:  payload[1],
		Amount:    amount,
		Change:    change,
		BHash:     payload[4],
		Timestamp: timestamp,
		DS:        ds,
		Payload:   payloadRaw}
	return
}

func (t tx) String() string {
	return t.Payload
}

func (m *Miner) SendChain() {
	m.chainMut.Lock()
	chain := make([]string, len(m.chain))
	for _, node := range m.chain {
		chain = append(chain, node.Pack())
	}
	m.chainMut.Unlock()
	utils.Write(strings.Join(chain, ";"), m.client)
}
