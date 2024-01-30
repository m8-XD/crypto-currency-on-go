package cryptography_test

import (
	"blockchain/pkg/cryptography"
	"strings"
	"testing"
)

func TestKeyGeneration(t *testing.T) {
	keyPair, err := cryptography.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := keyPair.Private()
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := keyPair.Public()
	if err != nil {
		t.Fatal(err)
	}
	newKeyPair, err := cryptography.GenerateKeyPairFromPrivate(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	newPrivateKey, err := newKeyPair.Private()
	if err != nil {
		t.Fatal(err)
	}
	newPublicKey, err := newKeyPair.Public()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(privateKey, newPrivateKey) {
		t.Fatal("private keys aren't equal")
	}
	if !strings.EqualFold(publicKey, newPublicKey) {
		t.Fatal("public keys aren't equal")
	}
}

func TestSigning(t *testing.T) {
	keyPair, err := cryptography.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	toSign := "foobar"

	privateKey, err := keyPair.Private()
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := keyPair.Public()
	if err != nil {
		t.Fatal(err)
	}
	hash, err := cryptography.HashMD5(toSign)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := cryptography.Sign(privateKey, hash)
	if !cryptography.Validate(publicKey, hash, signature) {
		t.Fatal("validation failed (valid signature)")
	}
	if cryptography.Validate(publicKey, []byte("wrong hash"), signature) {
		t.Fatal("validation failed (invalid signature)")
	}
}

func TestGenPublicWithRawPrivate(t *testing.T) {

	keyPair, err := cryptography.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := keyPair.Private()
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := keyPair.Public()
	if err != nil {
		t.Fatal(err)
	}

	newKeyPair, err := cryptography.GenerateKeyPairFromPrivate(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	newPrivateKey, err := newKeyPair.Private()
	if err != nil {
		t.Fatal(err)
	}
	newPublicKey, err := newKeyPair.Public()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(privateKey, newPrivateKey) {
		t.Fatal("private keys aren't equal")
	}
	if !strings.EqualFold(publicKey, newPublicKey) {
		t.Fatal("public keys aren't equal")
	}
}
