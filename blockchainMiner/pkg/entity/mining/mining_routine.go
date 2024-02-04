package mining

import (
	"blockchain/pkg/cryptography"
	"fmt"
	"strings"
)

// [+] validate ds
// [+] check if tx exists (input)
// [+]check if theres enough money in it
// mine

//one node = one tx
//check so tx cannot be before last node in chain

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
	m.chain = append(m.chain, node{})
	fmt.Println("niceeee")
	m.chainMut.Unlock()
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
