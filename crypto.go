package hedera

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Ed25519PrivateKey struct {
	keyData     []byte
	asStringRaw string
	chainCode   []byte
	publicKey   Ed25519PublicKey
}

func NewEd25519PrivateKey() Ed25519PrivateKey {
	_, privateKey := GenerateKeys()
	return privateKey
}

type Ed25519PublicKey struct {
	keyData     []byte
	asStringRaw string
	chainCode   []byte
}

func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(privateKey, publicKey)
	return publicKey, privateKey
}

func FromBytes(key []byte) string {
	s := hex.EncodeToString(key)
	fmt.Println(s)
	return s
}

func FromString(str string) []byte {
	b, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	return b
}
