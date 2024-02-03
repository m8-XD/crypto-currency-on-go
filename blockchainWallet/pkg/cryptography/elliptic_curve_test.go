package cryptography_test

import (
	"blockchain/pkg/cryptography"
	"bytes"
	"strings"
	"testing"
)

func TestKeyGeneration(t *testing.T) {
	keyPair, err := cryptography.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	privateKey := keyPair.PrivateHex()
	publicKey := keyPair.Public()
	newKeyPair, err := cryptography.GenKeysFromPrivate(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	newPrivateKey := newKeyPair.PrivateHex()
	newPublicKey := newKeyPair.Public()
	if !strings.EqualFold(privateKey, newPrivateKey) {
		t.Fatal("private keys aren't equal")
	}
	if !strings.EqualFold(string(publicKey), string(newPublicKey)) {
		t.Fatal("public keys aren't equal")
	}
}

func TestRecoveringKey(t *testing.T) {
	keyPair, err := cryptography.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	privateKey := keyPair.PrivateHex()
	publicKey := keyPair.Public()

	toSign := "foobar"

	signature, err := cryptography.Sign(privateKey, []byte(toSign))
	if err != nil {
		t.Fatal(err)
	}
	recoveredKey, err := cryptography.Recover([]byte(toSign), signature)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(publicKey, recoveredKey) {
		t.Fatal("keys are not equal!")
	}
}
