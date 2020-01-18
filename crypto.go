package hedera

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

const ed25519PrivateKeyPrefix = "302e020100300506032b657004220420"
const ed25519PubKeyPrefix = "302a300506032b6570032100"

type PublicKey interface {
	toProto() *proto.Key
}

func publicKeyFromProto(pbKey *proto.Key) (PublicKey, error) {
	switch key := pbKey.GetKey().(type) {
	case *proto.Key_Ed25519:
		return Ed25519PublicKeyFromBytes(key.Ed25519)

	case *proto.Key_ThresholdKey:
		threshold := key.ThresholdKey.GetThreshold()
		keys, err := publicKeyListFromProto(key.ThresholdKey.GetKeys())
		if err != nil {
			return nil, err
		}

		return NewThresholdKey(threshold).AddAll(keys), nil

	case *proto.Key_KeyList:
		keys, err := publicKeyListFromProto(key.KeyList)
		if err != nil {
			return nil, err
		}

		return NewKeyList().AddAll(keys), nil

	default:
		return nil, fmt.Errorf("key type not implemented: %v", key)
	}
}

func publicKeyListFromProto(pb *proto.KeyList) ([]PublicKey, error) {
	var keys []PublicKey = make([]PublicKey, len(pb.Keys))

	for i, pbKey := range pb.Keys {
		key, err := publicKeyFromProto(pbKey)

		if err != nil {
			return nil, err
		}

		keys[i] = key
	}

	return keys, nil
}

type Ed25519PrivateKey struct {
	keyData   []byte
	chainCode []byte
}

type Ed25519PublicKey struct {
	keyData []byte
}

func GenerateEd25519PrivateKey() (Ed25519PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	return Ed25519PrivateKey{
		keyData: privateKey,
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

	return Ed25519PrivateKey{
		keyData: privateKey,
	}, nil
}

func Ed25519PrivateKeyFromMnemonic(mnemonic Mnemonic, passPhrase string) (Ed25519PrivateKey, error) {
	salt := []byte("mnemonic" + passPhrase)
	seed := pbkdf2.Key([]byte(mnemonic.String()), salt, 2048, 64, sha512.New)

	h := hmac.New(sha512.New, []byte("ed25519 seed"))

	_, err := h.Write(seed)
	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	digest := h.Sum(nil)

	keyBytes := digest[0:32]
	chainCode := digest[32:len(digest)]

	// note the index is for derivation, not the index of the slice
	for _, index := range []uint32{44, 3030, 0, 0} {
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

func Ed25519PrivateKeyFromKeystore(ks []byte, passphrase string) (Ed25519PrivateKey, error) {
	return parseKeystore(ks, passphrase)
}

func Ed25519PrivateKeyReadKeystore(source io.Reader, passphrase string) (Ed25519PrivateKey, error) {
	keystoreBytes, err := ioutil.ReadAll(source)
	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	return Ed25519PrivateKeyFromKeystore(keystoreBytes, passphrase)
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

func Ed25519PublicKeyFromBytes(bytes []byte) (Ed25519PublicKey, error) {
	if len(bytes) != ed25519.PublicKeySize {
		return Ed25519PublicKey{}, fmt.Errorf("invalid public key")
	}

	return Ed25519PublicKey{
		keyData: bytes,
	}, nil
}

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

func (sk Ed25519PrivateKey) PublicKey() Ed25519PublicKey {
	return Ed25519PublicKey{
		keyData: sk.keyData[32:],
	}
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

func (sk Ed25519PrivateKey) Keystore(passphrase string) ([]byte, error) {
	return newKeystore(sk.keyData, passphrase)
}

func (sk Ed25519PrivateKey) WriteKeystore(destination io.Writer, passphrase string) error {
	keystore, err := sk.Keystore(passphrase)
	if err != nil {
		return err
	}

	_, err = destination.Write(keystore)

	return err
}

func (sk Ed25519PrivateKey) Sign(message []byte) []byte {
	return ed25519.Sign(sk.keyData, message)
}

func (sk Ed25519PrivateKey) SupportsDerivation() bool {
	return sk.chainCode != nil
}

// Derive a child key compatible with the iOS and Android wallets
// using a provided wallet/account index
//
// Use index 0 for the default account.
func (sk Ed25519PrivateKey) Derive(index uint32) (Ed25519PrivateKey, error) {
	if !sk.SupportsDerivation() {
		return Ed25519PrivateKey{}, fmt.Errorf("this private key does not support derivation")
	}

	derivedKeyBytes, chainCode := deriveChildKey(sk.Bytes(), sk.chainCode, index)

	derivedKey, err := Ed25519PrivateKeyFromBytes(derivedKeyBytes)

	if err != nil {
		return Ed25519PrivateKey{}, err
	}

	derivedKey.chainCode = chainCode

	return derivedKey, nil
}

func (pk Ed25519PublicKey) Bytes() []byte {
	return pk.keyData
}

func (pk Ed25519PublicKey) toProto() *proto.Key {
	return &proto.Key{Key: &proto.Key_Ed25519{Ed25519: pk.keyData}}
}
