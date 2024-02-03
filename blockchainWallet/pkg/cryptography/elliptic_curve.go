package cryptography

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyPair struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (kp *KeyPair) Public() []byte {
	publicKeyBytes := crypto.FromECDSAPub(kp.publicKey)
	return publicKeyBytes
}

func (kp *KeyPair) Private() []byte {
	privateKeyBytes := crypto.FromECDSA(kp.privateKey)
	return privateKeyBytes
}

func (kp *KeyPair) PrivateHex() string {
	privateKeyBytes := crypto.FromECDSA(kp.privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)
	return privateKeyHex
}

func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey
	return &KeyPair{privateKey, &publicKey}, nil
}

func GenKeysFromPrivate(privateKeySt string) (*KeyPair, error) {
	privateKeySt = strings.TrimPrefix(privateKeySt, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeySt)
	if err != nil {
		fmt.Println("!!!invalid private key!!!")
		return nil, err
	}
	return &KeyPair{privateKey, &privateKey.PublicKey}, nil
}

func Sign(privateKeyHex string, data []byte) (signature []byte, err error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(data)
	signature, err = crypto.Sign(hash.Bytes(), privateKey)
	return
}

// Recover key from digital signature
func Recover(data []byte, signature []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return nil, err
	}
	return sigPublicKey, nil
}
