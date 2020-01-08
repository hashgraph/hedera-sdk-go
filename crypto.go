package hedera

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/pbkdf2"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

const ed25519PrivateKeyPrefix = "302e020100300506032b657004220420"
const ed25519PubKeyPrefix = "302a300506032b6570032100"

var ErrNoRNG = errors.New("could not retrieve random bytes from the operating system")
var ErrUnderivable = errors.New("this private key does not support derivation")

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
		privateKey = ed25519.NewKeyFromSeed(bytes[0:32])

	default:
		return Ed25519PrivateKey{}, fmt.Errorf("invalid private key")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return Ed25519PrivateKey{
		keyData:   privateKey,
		publicKey: Ed25519PublicKey{publicKey},
	}, nil
}

func Ed25519PrivateKeyFromMnemonic(mnemonic string, passPhrase string) (Ed25519PrivateKey, error) {
	salt := []byte("mnemonic" + passPhrase)
	seed := pbkdf2.Key([]byte(mnemonic), salt, 2048, 64, sha512.New)

	h := hmac.New(sha512.New, []byte("ed25519 seed"))
	h.Write(seed)
	digest := h.Sum(nil)

	keyBytes := digest[0:32]
	chainCode := digest[32:len(digest)]

	// note the index is for derivation, not the index of the slice
	for _, index := range []uint32{ 44, 3030, 0, 0 } {
		keyBytes, chainCode = deriveChildKey(keyBytes, chainCode, index)
	}

	privateKey, err := Ed25519PrivateKeyFromBytes(keyBytes)

	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	privateKey.chainCode = chainCode

	return privateKey, nil
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

// todo: Ed25519PublicKeyFromBytes

// SLIP-10/BIP-32 Child Key derivation
func deriveChildKey(parentKey []byte, chainCode []byte, index uint32) ([]byte, []byte) {
	h := hmac.New(sha512.New, chainCode)

	input := make([]byte, 37)

	// 0x00 + parentKey + index(BE)
	input[0] = 0

	copy(input[1:37], parentKey)

	binary.BigEndian.PutUint32(input[33:37], index)

	// harden the input
	input[33] |= 128

	h.Write(input)
	digest := h.Sum(nil)

	return digest[0:32], digest[32:len(digest)]
}

type MnemonicResult struct {
	mnemonic string
}

// todo: rename as toPrivateKey
func (mr MnemonicResult) GenerateKey(passPhrase string) (Ed25519PrivateKey, error) {
	return Ed25519PrivateKeyFromMnemonic(mr.mnemonic, passPhrase)
}

// Generate a random 24-word mnemonic
func GenerateMnemonic() (*MnemonicResult, error) {
	entropy, err := bip39.NewEntropy(256)

	if err != nil {
		// It is only possible for there to be an error if the operating
		// system's rng is unreadable
		return nil, ErrNoRNG
	}

	mnemonic, err := bip39.NewMnemonic(entropy)

	if err != nil {
		// todo: return proper error
		return nil, err
	}

	return &MnemonicResult{mnemonic}, nil
}

// todo: Mnemonic From String

// todo: Mnemonic To String

// todo: func NewMnemonic([]string)

// todo: Words -> []string

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

func (sk Ed25519PrivateKey) SupportsDerivation() bool {
	return sk.chainCode != nil
}

// Given a wallet/account index, derive a child key compatible with the iOS and Android wallets.
//
// Use index 0 for the default account.
func (sk Ed25519PrivateKey) Derive(index uint32) (Ed25519PrivateKey, error) {
	if !sk.SupportsDerivation() {
		return Ed25519PrivateKey{}, ErrUnderivable
	}

	derivedKeyBytes, chainCode := deriveChildKey(sk.Bytes(), sk.chainCode, index)

	derivedKey, err := Ed25519PrivateKeyFromBytes(derivedKeyBytes)

	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	derivedKey.chainCode = chainCode

	return derivedKey, nil
}
