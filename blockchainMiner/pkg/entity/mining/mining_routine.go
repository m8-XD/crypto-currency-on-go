package mining

import (
	"blockchain/pkg/cryptography"
	"blockchain/pkg/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// one node = one tx
// check so tx cannot be before last node in chain
const miningDifficulty = 5

func Mine(m *Miner, tx tx) {

	ok, err := validate(tx)
	if !ok || err != nil {
		fmt.Println("digital signature validation failed")
		return
	}
	blockInd, ok := referedNodeInd(m, tx)
	if !ok {
		fmt.Println("TX doesn't exists")
		return
	}
	if !enoughMoney(m, tx, blockInd) {
		fmt.Println("not enough money in current TX")
		return
	}

	m.chainMut.Lock()
	if !checkTimeValidity(m, tx) {
		fmt.Println("current tx was created before latest, skipping...")
		return
	}
	n := mineNode(m, tx)

	m.chain = append(m.chain, n)
	m.chainMut.Unlock()

	utils.Write(n.Pack(), m.client)
}

func mineNode(m *Miner, tx tx) node {
	prefix := strings.Repeat("0", miningDifficulty)

	lastNode := m.chain[len(m.chain)-1]

	pHead := lastNode.Header
	timestamp := tx.Timestamp

	payload := pHead + tx.String() + fmt.Sprint(timestamp)

	header := ""
	headerBytes := make([]byte, 0)
	var nonce int64 = 0
	for !strings.HasPrefix(header, prefix) {
		nonce++
		hash := sha256.Sum256([]byte(payload + fmt.Sprint(nonce)))
		headerBytes = hash[:]
		header = hex.EncodeToString(headerBytes)
	}
	return node{
		Header:    header,
		PHeader:   pHead,
		Nonce:     nonce,
		TX:        tx,
		Timestamp: timestamp,
	}
}

func validate(tx tx) (bool, error) {
	payload := tx.Payload
	restoredKey, err := cryptography.Recover([]byte(payload), []byte(tx.DS))
	if err != nil {
		return false, err
	}
	restoredWAddr := cryptography.WaletAddr(restoredKey)
	return strings.EqualFold(restoredWAddr, tx.WAddr), nil
}

func referedNodeInd(m *Miner, tx tx) (int, bool) {
	for i := 0; i < len(m.chain); i++ {
		if strings.EqualFold(m.chain[i].Header, tx.BHash) {
			return i, true
		}
	}
	return -1, false
}

func enoughMoney(m *Miner, tx tx, blockInd int) bool {
	referedTX := m.chain[blockInd].TX
	if !strings.EqualFold(referedTX.RecWAddr, tx.WAddr) {
		return checkChange(referedTX, tx)
	}
	if referedTX.Amount >= tx.Amount+tx.Change {
		return true
	}
	return false
}

func checkChange(refered tx, current tx) bool {
	if !strings.EqualFold(refered.WAddr, current.WAddr) {
		return false
	}
	return refered.Change >= current.Amount+current.Change
}

func checkTimeValidity(m *Miner, tx tx) bool {
	return m.chain[len(m.chain)-1].Timestamp < tx.Timestamp
}
