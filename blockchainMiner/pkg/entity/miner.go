package entity

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Miner struct {
	client   *Client
	txs      []tx //transactions
	chain    []node
	txMut    sync.Mutex
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
}

func (m *Miner) IsRunning() bool {
	return m.client.isRunning
}

func (m *Miner) Client() *Client {
	return m.client
}

func (m *Miner) SetClient(c *Client) {
	m.client = c
}

func (m *Miner) Start() {
	m.client.Start()
	m.txs = make([]tx, 0)
	m.chain = make([]node, 0)
}

func (m *Miner) AddTX(txRaw string) {
	tx, err := parseTX(txRaw)
	if err == nil {
		m.txMut.Lock()
		m.txs = append(m.txs, tx)
		m.txMut.Unlock()
	}
}

func (m *Miner) NextTx() (*tx, bool) {
	if len(m.txs) == 0 {
		return nil, false
	}
	last := m.txs[len(m.txs)-1]
	m.txs = m.txs[:len(m.txs)-1]
	return &last, true
}

func parseTX(txRaw string) (txn tx, err error) {
	parsedTX := strings.Split(txRaw, ":")
	if len(parsedTX) != 2 {
		err = errors.New("invalid TX")
		return
	}
	payloadRaw, ds := parsedTX[0], parsedTX[1]
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
		DS:        ds}
	return
}