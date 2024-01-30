package cryptography

import (
	"blockchain/pkg/ui"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

const privateKeyPref = "-----BEGIN EC PRIVATE KEY-----\n"
const privateKeySuff = "-----END EC PRIVATE KEY-----\n"
const publicKeyPref = "-----BEGIN PUBLIC KEY-----\n"
const publicKeySuff = "-----END PUBLIC KEY-----\n"

type KeyPair struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (kp *KeyPair) Private() (string, error) {

	encoded, err := x509.MarshalECPrivateKey(kp.privateKey)

	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})
	key := trimPrivateKey(string(pemEncoded))
	return key, nil
}

func (kp *KeyPair) Public() (string, error) {

	encoded, err := x509.MarshalPKIXPublicKey(kp.publicKey)

	if err != nil {
		return "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	key := trimPublicKey(string(pemEncodedPub))
	return key, nil
}

func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		fmt.Println("cannot generate key pair: " + err.Error())
		return nil, err
	}

	publicKey := generatePublicKey(privateKey)
	return &KeyPair{privateKey, publicKey}, nil
}

func GenerateKeyPairFromPrivate(privateKeySt string) (*KeyPair, error) {
	privateKey, err := decodePrivate(privateKeySt)

	if err != nil {
		ui.Error("invalid private key")
		fmt.Println(err)
		return nil, err
	}
	return &KeyPair{privateKey, &privateKey.PublicKey}, nil
}

func Sign(privateKeySt string, hash []byte) (signature []byte, err error) {
	privateKey, err := decodePrivate(privateKeySt)
	if err != nil {
		return nil, err
	}
	signature, err = ecdsa.SignASN1(rand.Reader, privateKey, hash)
	return
}

func Validate(publicKeySt string, hash []byte, signature []byte) bool {
	publicKey, err := decodePublic(publicKeySt)
	if err != nil {
		fmt.Println("invalid public key to validate")
		return false
	}
	return ecdsa.VerifyASN1(publicKey, hash, signature)
}

func wrapPrivate(key string) string {
	key = strings.Replace(key, "@", "\n", -1)
	key = privateKeyPref + key + privateKeySuff
	return key
}

func wrapPublic(key string) string {
	key = strings.Replace(key, "@", "\n", -1)
	key = publicKeyPref + key + publicKeySuff
	return key
}
func generatePublicKey(privateKey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	return &privateKey.PublicKey
}

func decodePrivate(pemEncodedPriv string) (privateKey *ecdsa.PrivateKey, err error) {
	pemEncodedPriv = wrapPrivate(pemEncodedPriv)
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))

	x509EncodedPriv := blockPriv.Bytes

	privateKey, err = x509.ParseECPrivateKey(x509EncodedPriv)

	return
}

func decodePublic(pemEncodedPub string) (publicKey *ecdsa.PublicKey, err error) {
	pemEncodedPub = wrapPublic(pemEncodedPub)
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey = genericPublicKey.(*ecdsa.PublicKey)

	return
}

func trimPrivateKey(privateKey string) string {
	privateKey = strings.TrimPrefix(privateKey, privateKeyPref)
	privateKey = strings.TrimSuffix(privateKey, privateKeySuff)
	privateKey = strings.Replace(privateKey, "\n", "@", -1)
	return privateKey
}

func trimPublicKey(publicKey string) string {
	publicKey = strings.TrimPrefix(publicKey, publicKeyPref)
	publicKey = strings.TrimSuffix(publicKey, publicKeySuff)
	publicKey = strings.Replace(publicKey, "\n", "@", -1)
	return publicKey
}
