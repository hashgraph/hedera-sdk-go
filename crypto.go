package hedera

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

const ed25519PrivateKeyPrefix = "302e020100300506032b657004220420"
const ed25519PubKeyPrefix = "302a300506032b6570032100"

type Ed25519PrivateKey struct {
	keyData   []byte
	chainCode []byte
	publicKey Ed25519PublicKey
}

type Ed25519PublicKey struct {
	keyData []byte
}

func GenerateEd25519PrivateKey() (Ed25519PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Ed25519PrivateKey{}, nil
	}

	return Ed25519PrivateKey{
		keyData:   privateKey,
		publicKey: Ed25519PublicKey{publicKey},
	}, nil
}

func Ed25519PrivateKeyFromBytes(bytes []byte) (Ed25519PrivateKey, error) {
	var privateKey ed25519.PrivateKey

	switch len(bytes) {
	case 32:
		// The bytes array has just the private key
		privateKey = ed25519.NewKeyFromSeed(bytes)

	case 64:
		privateKey = bytes

	default:
		return Ed25519PrivateKey{}, fmt.Errorf("invalid private key")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return Ed25519PrivateKey{
		keyData:   privateKey,
		publicKey: Ed25519PublicKey{publicKey},
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
		if strings.HasPrefix(s, ed25519PrivateKeyPrefix) {
			return Ed25519PrivateKeyFromString(s[32:])
		}
	}

	return Ed25519PrivateKey{}, fmt.Errorf("invalid private key with length %v", len(s))
}

func Ed25519PublicKeyFromString(s string) (Ed25519PublicKey, error) {
	switch len(s) {
	case 64: // raw public key
		bytes, err := hex.DecodeString(s)
		if err != nil {
			return Ed25519PublicKey{}, err
		}

		return Ed25519PublicKey{bytes}, nil

	case 88: // DER encoded public key
		if strings.HasPrefix(s, ed25519PubKeyPrefix) {
			pk, err := Ed25519PublicKeyFromString(s[24:])
			if err != nil {
				return Ed25519PublicKey{}, err
			}
			return pk, nil
		}
	}
	return Ed25519PublicKey{}, fmt.Errorf("invalid public key with length %v", len(s))
}

func (sk Ed25519PrivateKey) PublicKey() Ed25519PublicKey {
	return sk.publicKey
}

func (sk Ed25519PrivateKey) String() string {
	return fmt.Sprint(ed25519PrivateKeyPrefix, hex.EncodeToString(sk.keyData[:32]))
}

func (pk Ed25519PublicKey) String() string {
	return fmt.Sprint(ed25519PubKeyPrefix, hex.EncodeToString(pk.keyData))
}

func (sk Ed25519PrivateKey) Bytes() []byte {
	return sk.keyData
}

func (pk Ed25519PublicKey) Bytes() []byte {
	return pk.keyData
}

func (pk Ed25519PublicKey) toProto() *proto.Key {
	return &proto.Key{Key: &proto.Key_Ed25519{Ed25519: pk.keyData}}
}

func (sk Ed25519PrivateKey) Sign(message []byte) []byte {
	return ed25519.Sign(sk.keyData, message)
}
