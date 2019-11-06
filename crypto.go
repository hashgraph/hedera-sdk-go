package hedera

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

const ed25519PrivKeyPrefix = "302e020100300506032b657004220420"
const ed25519PubKeyPrefix = "302a300506032b6570032100"

type Ed25519PrivateKey struct {
	keyData   []byte
	chainCode []byte
	publicKey Ed25519PublicKey
}

func GenerateEd25519PrivateKey() (Ed25519PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Ed25519PrivateKey{}, nil
	}

	return Ed25519PrivateKey{
		keyData: privateKey,
		publicKey: Ed25519PublicKey{
			keyData: publicKey,
		},
	}, nil
}

type Ed25519PublicKey struct {
	keyData []byte
}

func Ed25519PrivateKeyFromBytes(bytes []byte) (Ed25519PrivateKey, error) {
	var privateKey ed25519.PrivateKey

	switch len(bytes) {
	case 32:
		// The bytes array has just the private key
		privateKey = ed25519.NewKeyFromSeed(bytes)

	case 64:
		privateKey = ed25519.PrivateKey(bytes)

	default:
		return Ed25519PrivateKey{}, fmt.Errorf("invalid private key")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return Ed25519PrivateKey{
		keyData: privateKey,
		publicKey: Ed25519PublicKey{
			keyData: publicKey,
		},
	}, nil
}

func Ed25519PrivateKeyFromString(s string) (Ed25519PrivateKey, error) {
	switch len(s) {
	case 64, 128: // private key : public key
		bytes, err := hex.DecodeString(s)
		if err != nil {
			return Ed25519PrivateKey{}, err
		}

		return Ed25519PrivateKeyFromBytes(bytes)

	case 96: // prefix-encoded private key
		if strings.HasPrefix(s, ed25519PrivKeyPrefix) {
			return Ed25519PrivateKeyFromString(s[32:])
		}
	}

	return Ed25519PrivateKey{}, fmt.Errorf("invalid private key with length %v", len(s))
}

func (priv Ed25519PrivateKey) PublicKey() Ed25519PublicKey {
	return priv.publicKey
}

func (priv Ed25519PrivateKey) String() string {
	return fmt.Sprint(ed25519PrivKeyPrefix, hex.EncodeToString(priv.keyData[:32]))
}

func (pub Ed25519PublicKey) String() string {
	return fmt.Sprint(ed25519PubKeyPrefix, hex.EncodeToString(pub.keyData))
}
