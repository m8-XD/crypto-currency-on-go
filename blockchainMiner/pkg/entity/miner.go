package entity

type Miner struct {
	client *Client
	txs    []tx //transactions
	chain  []node
}

// transaction
type tx struct {
	wAddr     string
	recWAddr  string //reciever wallet address
	amount    float64
	change    float64
	timestamp int64
	ds        string //digital signature
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
	m.txs = append(m.txs, parseTX(txRaw))
}

func parseTX(txRaw string) tx {
}
