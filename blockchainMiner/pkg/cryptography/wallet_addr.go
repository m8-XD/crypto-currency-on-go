package cryptography

import (
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func WaletAddr(pubKey []byte) (wAddr string) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey[1:])
	wAddr = hexutil.Encode(hash.Sum(nil)[12:])

	return
}
