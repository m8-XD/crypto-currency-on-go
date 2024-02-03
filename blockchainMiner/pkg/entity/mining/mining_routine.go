package mining

import (
	"blockchain/pkg/cryptography"
	"fmt"
	"strings"
)

// validate ds
// check if tx exists (input)
// check if theres enough money in it
// mine

//one node = one tx
//check so tx cannot be before last node in chain

func Mine(m *Miner, tx tx) {

	ok, err := validate(tx)
	if !ok || err != nil {
		fmt.Println("digital signature validation failed")
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
