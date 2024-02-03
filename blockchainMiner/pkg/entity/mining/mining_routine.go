package mining

// validate ds 
// check if tx exists (input)
// check if theres enough money in it
// mine

func Mine(m *Miner) {
	m.chainMut.Lock()
    implement mining routine
	m.chain = append(m.chain, node{})
	m.LastBlockTs = m.LastTXts
	m.chainMut.Unlock()
}
