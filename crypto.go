package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	ecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	protobuf "google.golang.org/protobuf/proto"
)

const _Ed25519PrivateKeyPrefix = "302e020100300506032b657004220420"

type Key interface {
	_ToProtoKey() *services.Key
	String() string
}

func KeyFromBytes(bytes []byte) (Key, error) {
	protoKey := &services.Key{}

	err := protobuf.Unmarshal(bytes, protoKey)
	if err != nil {
		return nil, err
	}

	return _KeyFromProtobuf(protoKey)
}

func KeyToBytes(key Key) ([]byte, error) {
	protoKey := key._ToProtoKey()
	return protobuf.Marshal(protoKey)
}

func _KeyFromProtobuf(pbKey *services.Key) (Key, error) {
	if pbKey == nil {
		return PublicKey{}, errParameterNull
	}
	switch key := pbKey.GetKey().(type) {
	case *services.Key_Ed25519:
		return PublicKeyFromBytesEd25519(key.Ed25519)

	case *services.Key_ThresholdKey:
		threshold := int(key.ThresholdKey.GetThreshold())
		keys, err := _KeyListFromProtobuf(key.ThresholdKey.GetKeys())
		if err != nil {
			return nil, err
		}
		keys.threshold = threshold

		return &keys, nil

	case *services.Key_KeyList:
		keys, err := _KeyListFromProtobuf(key.KeyList)
		if err != nil {
			return nil, err
		}

		return &keys, nil

	case *services.Key_ContractID:
		return _ContractIDFromProtobuf(key.ContractID), nil

	case *services.Key_ECDSASecp256K1:
		return PublicKeyFromBytesECDSA(key.ECDSASecp256K1)

	case *services.Key_DelegatableContractId:
		return _DelegatableContractIDFromProtobuf(key.DelegatableContractId), nil

	default:
		return nil, _NewErrBadKeyf("key type not implemented: %v", key)
	}
}

type PrivateKey struct {
	ecdsaPrivateKey   *_ECDSAPrivateKey
	ed25519PrivateKey *_Ed25519PrivateKey
}

type PublicKey struct {
	ecdsaPublicKey   *_ECDSAPublicKey
	ed25519PublicKey *_Ed25519PublicKey
}

/**
 *  SDK needs to provide  a way to set an unusable key such as an Ed25519 all-zeros
 *  key, since it is (presumably) impossible to find the 32-byte string whose SHA-512 hash begins with 32 bytes
 *  of zeros. We recommend using all-zeros to clearly advertise any unusable keys.
 */
func ZeroKey() (PublicKey, error) {
	return PublicKeyFromString("0000000000000000000000000000000000000000000000000000000000000000")
}

// PrivateKeyGenerateEcdsa Generates a new ECDSASecp256K1 key
func PrivateKeyGenerateEcdsa() (PrivateKey, error) {
	key, err := _GenerateECDSAPrivateKey()
	if err != nil {
		return PrivateKey{}, err
	}
	return PrivateKey{
		ecdsaPrivateKey: key,
	}, nil
}

// Deprecated: use `PrivateKeyGenerateEd25519()` instead
func PrivateKeyGenerate() (PrivateKey, error) {
	return PrivateKeyGenerateEd25519()
}

// PrivateKeyGenerateEd25519 Generates a new Ed25519 key
func PrivateKeyGenerateEd25519() (PrivateKey, error) {
	key, err := _GenerateEd25519PrivateKey()
	if err != nil {
		return PrivateKey{}, err
	}
	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

// Deprecated the use of raw bytes for a Ed25519 private key is deprecated; use PrivateKeyFromBytesEd25519() instead.
func PrivateKeyFromBytes(bytes []byte) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromBytes(bytes)
	if err != nil {
		key2, err2 := _ECDSAPrivateKeyFromBytes(bytes)
		if err2 != nil {
			return PrivateKey{}, err2
		}

		return PrivateKey{
			ecdsaPrivateKey: key2,
		}, nil
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyFromBytesDer(bytes []byte) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromBytes(bytes)
	if err != nil {
		key2, err2 := _ECDSAPrivateKeyFromBytes(bytes)
		if err2 != nil {
			return PrivateKey{}, err2
		}

		return PrivateKey{
			ecdsaPrivateKey: key2,
		}, nil
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyFromBytesEd25519(bytes []byte) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromBytes(bytes)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyFromBytesECDSA(bytes []byte) (PrivateKey, error) {
	key, err := _ECDSAPrivateKeyFromBytes(bytes)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ecdsaPrivateKey: key,
	}, nil
}

func PublicKeyFromBytesEd25519(bytes []byte) (PublicKey, error) {
	key, err := _Ed25519PublicKeyFromBytes(bytes)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		ed25519PublicKey: key,
	}, nil
}

func PublicKeyFromBytesECDSA(bytes []byte) (PublicKey, error) {
	key, err := _ECDSAPublicKeyFromBytes(bytes)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		ecdsaPublicKey: key,
	}, nil
}

// Deprecated the use of raw bytes for a Ed25519 private key is deprecated; use PublicKeyFromBytesEd25519() instead.
func PublicKeyFromBytes(bytes []byte) (PublicKey, error) {
	key, err := _Ed25519PublicKeyFromBytes(bytes)
	if err != nil {
		key2, err2 := _ECDSAPublicKeyFromBytes(bytes)
		if err2 != nil {
			return PublicKey{}, err2
		}

		return PublicKey{
			ecdsaPublicKey: key2,
		}, nil
	}

	return PublicKey{
		ed25519PublicKey: key,
	}, nil
}

func PublicKeyFromBytesDer(bytes []byte) (PublicKey, error) {
	key, err := _Ed25519PublicKeyFromBytes(bytes)
	if err != nil {
		key2, err2 := _ECDSAPublicKeyFromBytes(bytes)
		if err2 != nil {
			return PublicKey{}, err2
		}

		return PublicKey{
			ecdsaPublicKey: key2,
		}, nil
	}

	return PublicKey{
		ed25519PublicKey: key,
	}, nil
}

// Deprecated
// PrivateKeyFromMnemonic recovers an _Ed25519PrivateKey from a valid 24 word length mnemonic phrase and a
// passphrase.
//
// An empty string can be passed for passPhrase If the mnemonic phrase wasn't generated with a passphrase. This is
// required to recover a private key from a mnemonic generated by the Android and iOS wallets.
func PrivateKeyFromMnemonic(mnemonic Mnemonic, passPhrase string) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromMnemonic(mnemonic, passPhrase)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

// The use of raw bytes for a Ed25519 private key is deprecated; use PrivateKeyFromStringEd25519() instead.
func PrivateKeyFromString(s string) (PrivateKey, error) {
	byt, err := hex.DecodeString(s)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKeyFromBytes(byt)
}

// PrivateKeyFromStringDer Creates PrivateKey from hex string with a der prefix
func PrivateKeyFromStringDer(s string) (PrivateKey, error) {
	KeyEd25519, err := _Ed25519PrivateKeyFromString(s)
	if err == nil {
		return PrivateKey{ed25519PrivateKey: KeyEd25519}, nil
	}

	keyECDSA, err := _ECDSAPrivateKeyFromString(s)
	if err == nil {
		return PrivateKey{ecdsaPrivateKey: keyECDSA}, nil
	}

	return PrivateKey{}, errors.New("invalid private key format")
}

func PrivateKeyFromStringEd25519(s string) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromString(s)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

// Deprecated: use PrivateKeyFromStringECDSA() instead
func PrivateKeyFromStringECSDA(s string) (PrivateKey, error) {
	return PrivateKeyFromStringECDSA(s)
}

func PrivateKeyFromStringECDSA(s string) (PrivateKey, error) {
	trimmedKey := strings.TrimPrefix(s, "0x")
	key, err := _ECDSAPrivateKeyFromString(trimmedKey)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ecdsaPrivateKey: key,
	}, nil
}

func PrivateKeyFromSeedEd25519(seed []byte) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromSeed(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyFromSeedECDSAsecp256k1(seed []byte) (PrivateKey, error) {
	key, err := _ECDSAPrivateKeyFromSeed(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ecdsaPrivateKey: key,
	}, nil
}

func PrivateKeyFromKeystore(ks []byte, passphrase string) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromKeystore(ks, passphrase)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

// PrivateKeyReadKeystore recovers an _Ed25519PrivateKey from an encrypted _Keystore file.
func PrivateKeyReadKeystore(source io.Reader, passphrase string) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyReadKeystore(source, passphrase)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyFromPem(bytes []byte, passphrase string) (PrivateKey, error) {
	key, err := _Ed25519PrivateKeyFromPem(bytes, passphrase)
	if err != nil {
		key, err := _ECDSAPrivateKeyFromPem(bytes, passphrase)
		if err != nil {
			return PrivateKey{}, err
		}
		return PrivateKey{
			ecdsaPrivateKey: key,
		}, nil
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

func PrivateKeyReadPem(source io.Reader, passphrase string) (PrivateKey, error) {
	// note: Passphrases are currently not supported, but included in the function definition to avoid breaking
	// changes in the future.

	key, err := _Ed25519PrivateKeyReadPem(source, passphrase)
	if err != nil {
		key, err := _ECDSAPrivateKeyReadPem(source, passphrase)
		if err != nil {
			return PrivateKey{}, err
		}
		return PrivateKey{
			ecdsaPrivateKey: key,
		}, nil
	}

	return PrivateKey{
		ed25519PrivateKey: key,
	}, nil
}

// The use of raw bytes for a Ed25519 public key is deprecated; use PublicKeyFromStringEd25519/ECDSA() instead.
func PublicKeyFromString(s string) (PublicKey, error) {
	byt, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKeyFromBytes(byt)
}

func PublicKeyFromStringECDSA(s string) (PublicKey, error) {
	key, err := _ECDSAPublicKeyFromString(s)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		ecdsaPublicKey: key,
	}, nil
}

func PublicKeyFromStringEd25519(s string) (PublicKey, error) {
	key, err := _Ed25519PublicKeyFromString(s)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		ed25519PublicKey: key,
	}, nil
}

func _DeriveEd25519ChildKey(parentKey []byte, chainCode []byte, index uint32) ([]byte, []byte, error) {
	if IsHardenedIndex(index) {
		return nil, nil, errors.New("the index should not be pre-hardened")
	}

	h := hmac.New(sha512.New, chainCode)

	input := make([]byte, 37)

	// 0x00 + parentKey + _Index(BE)
	input[0] = 0

	copy(input[1:37], parentKey)

	binary.BigEndian.PutUint32(input[33:37], index)

	// harden the input
	input[33] |= 128

	if _, err := h.Write(input); err != nil {
		return nil, nil, err
	}

	digest := h.Sum(nil)

	return digest[0:32], digest[32:], nil
}

func _DeriveECDSAChildKey(parentKey []byte, chainCode []byte, index uint32) ([]byte, []byte, error) {
	h := hmac.New(sha512.New, chainCode)

	isHardened := IsHardenedIndex(index)
	input := make([]byte, 37)
	if len(parentKey) != 32 {
		return nil, nil, fmt.Errorf("invalid private key length")
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(parentKey)

	if isHardened {
		offset := 33 - len(parentKey)
		copy(input[offset:], parentKey)
	} else {
		copy(input, pubKey.SerializeCompressed())
	}

	binary.BigEndian.PutUint32(input[33:37], index)

	if _, err := h.Write(input); err != nil {
		return nil, nil, err
	}

	i := h.Sum(nil)

	il := new(big.Int)
	il.SetBytes(i[0:32])
	ir := i[32:]

	ki := new(big.Int)
	ki.Add(privKey.ToECDSA().D, il)
	ki.Mod(ki, privKey.ToECDSA().Curve.Params().N)

	return ki.Bytes(), ir, nil
}

func _DeriveLegacyChildKey(parentKey []byte, index int64) ([]byte, error) {
	in := make([]uint8, 8)

	switch switchIndex := index; {
	case switchIndex == int64(0xffffffffff):
		in = []uint8{0, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff}
	case switchIndex > 0xffffffff:
		return nil, errors.New("derive index is out of range")
	default:
		if switchIndex < 0 {
			for i := 0; i < 4; i++ {
				in[i] = uint8(0xff)
			}
		}

		for i := 4; i < len(in); i++ {
			in[i] = uint8(switchIndex)
		}
	}

	password := make([]uint8, len(parentKey))
	copy(password, parentKey)
	password = append(password, in...)

	salt := []byte{0xFF}

	return pbkdf2.Key(password, salt, 2048, 32, sha512.New), nil
}

func (sk PrivateKey) PublicKey() PublicKey {
	if sk.ecdsaPrivateKey != nil {
		return PublicKey{
			ecdsaPublicKey: sk.ecdsaPrivateKey._PublicKey(),
		}
	}

	if sk.ed25519PrivateKey != nil {
		return PublicKey{
			ed25519PublicKey: sk.ed25519PrivateKey._PublicKey(),
		}
	}

	return PublicKey{}
}

func (sk PrivateKey) ToAccountID(shard uint64, realm uint64) *AccountID {
	return sk.PublicKey().ToAccountID(shard, realm)
}

func (pk PublicKey) ToAccountID(shard uint64, realm uint64) *AccountID {
	temp := pk

	return &AccountID{
		Shard:    shard,
		Realm:    realm,
		Account:  0,
		AliasKey: &temp,
		checksum: nil,
	}
}

// String returns the text-encoded representation of the PrivateKey.
func (sk PrivateKey) String() string {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._StringDer()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._StringDer()
	}

	return ""
}

func (sk PrivateKey) StringRaw() string {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._StringRaw()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._StringRaw()
	}

	return ""
}

func (sk PrivateKey) StringDer() string {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._StringDer()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._StringDer()
	}

	return ""
}

func (pk PublicKey) String() string {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._StringDer()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._StringDer()
	}

	return ""
}

func (pk PublicKey) StringDer() string {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._StringDer()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._StringDer()
	}

	return ""
}

func (pk PublicKey) StringRaw() string {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._StringRaw()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._StringRaw()
	}

	return ""
}

// `Deprecated: Use ToEvmAddress instead`
func (pk PublicKey) ToEthereumAddress() string {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._ToEthereumAddress()
	}

	panic("unsupported operation on Ed25519PublicKey")
}

func (pk PublicKey) ToEvmAddress() string {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._ToEthereumAddress()
	}

	panic("unsupported operation on Ed25519PublicKey")
}

/*
 * For `Ed25519` the result of this method call is identical to `toBytesRaw()` while for `ECDSA`
 * this method is identical to `toBytesDer()`.
 *
 * We strongly recommend using `toBytesRaw()` or `toBytesDer()` instead.
 */
func (sk PrivateKey) Bytes() []byte {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._BytesDer()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._BytesRaw()
	}

	return []byte{}
}

func (sk PrivateKey) BytesDer() []byte {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._BytesDer()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._BytesDer()
	}

	return []byte{}
}

func (sk PrivateKey) BytesRaw() []byte {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._BytesRaw()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._BytesRaw()
	}

	return []byte{}
}

/*
 * For `Ed25519` the result of this method call is identical to `toBytesRaw()` while for `ECDSA`
 * this method is identical to `toBytesDer()`.
 *
 * We strongly recommend using `toBytesRaw()` or `toBytesDer()` instead.
 */
func (pk PublicKey) Bytes() []byte {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._BytesDer()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._BytesRaw()
	}

	return []byte{}
}

func (pk PublicKey) BytesRaw() []byte {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._BytesRaw()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._BytesRaw()
	}

	return []byte{}
}

func (pk PublicKey) BytesDer() []byte {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._BytesDer()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._BytesDer()
	}

	return []byte{}
}

func (sk PrivateKey) Keystore(passphrase string) ([]byte, error) {
	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._Keystore(passphrase)
	}

	return []byte{}, errors.New("only ed25519 keystore is supported right now")
}

func (sk PrivateKey) WriteKeystore(destination io.Writer, passphrase string) error {
	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._WriteKeystore(destination, passphrase)
	}

	return errors.New("only writing ed25519 keystore is supported right now")
}

// Sign signs the provided message with the Ed25519PrivateKey.
func (sk PrivateKey) Sign(message []byte) []byte {
	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._Sign(message)
	}
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._Sign(message)
	}

	return []byte{}
}

func (sk PrivateKey) SupportsDerivation() bool {
	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._SupportsDerivation()
	}

	return false
}

func (sk PrivateKey) Derive(index uint32) (PrivateKey, error) {
	if sk.ed25519PrivateKey != nil {
		key, err := sk.ed25519PrivateKey._Derive(index)
		if err != nil {
			return PrivateKey{}, err
		}

		return PrivateKey{
			ed25519PrivateKey: key,
		}, nil
	}
	if sk.ecdsaPrivateKey != nil {
		key, err := sk.ecdsaPrivateKey._Derive(index)
		if err != nil {
			return PrivateKey{}, err
		}

		return PrivateKey{
			ecdsaPrivateKey: key,
		}, nil
	}

	return PrivateKey{}, nil
}

func (sk PrivateKey) LegacyDerive(index int64) (PrivateKey, error) {
	if sk.ed25519PrivateKey != nil {
		key, err := sk.ed25519PrivateKey._LegacyDerive(index)
		if err != nil {
			return PrivateKey{}, err
		}

		return PrivateKey{
			ed25519PrivateKey: key,
		}, nil
	}

	return PrivateKey{}, errors.New("only ed25519 legacy derivation is supported")
}

func (sk PrivateKey) _ToProtoKey() *services.Key {
	if sk.ecdsaPrivateKey != nil {
		return sk.ecdsaPrivateKey._ToProtoKey()
	}

	if sk.ed25519PrivateKey != nil {
		return sk.ed25519PrivateKey._ToProtoKey()
	}

	return &services.Key{}
}

func (pk PublicKey) _ToProtoKey() *services.Key {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._ToProtoKey()
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._ToProtoKey()
	}

	return &services.Key{}
}

func (pk PublicKey) _ToSignaturePairProtobuf(signature []byte) *services.SignaturePair {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._ToSignaturePairProtobuf(signature)
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._ToSignaturePairProtobuf(signature)
	}

	return &services.SignaturePair{}
}

// SignTransaction signes the transaction and adds the signature to the transaction
func (sk PrivateKey) SignTransaction(tx TransactionInterface) ([]byte, error) {
	baseTx := tx.getBaseTransaction()

	if sk.ecdsaPrivateKey != nil {
		b, err := sk.ecdsaPrivateKey._SignTransaction(baseTx)
		if err != nil {
			return []byte{}, err
		}

		return b, nil
	}

	if sk.ed25519PrivateKey != nil {
		b, err := sk.ed25519PrivateKey._SignTransaction(baseTx)
		if err != nil {
			return []byte{}, err
		}

		return b, nil
	}

	return []byte{}, errors.New("key type not supported, only ed25519 and ECDSASecp256K1 are supported right now")
}

func (pk PublicKey) Verify(message []byte, signature []byte) bool {
	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._Verify(message, signature)
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._Verify(message, signature)
	}

	return false
}

func (pk PublicKey) VerifyTransaction(tx TransactionInterface) bool {
	baseTx := tx.getBaseTransaction()

	if pk.ecdsaPublicKey != nil {
		return pk.ecdsaPublicKey._VerifyTransaction(baseTx)
	}

	if pk.ed25519PublicKey != nil {
		return pk.ed25519PublicKey._VerifyTransaction(baseTx)
	}

	return false
}

func Keccak256Hash(data []byte) (h Hash) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	copy(h[:], hash.Sum(nil))
	return h
}

func VerifySignature(pubkey, digestHash, signature []byte) bool {
	pubKey, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		return false
	}

	recoveredKey, _, err := ecdsa.RecoverCompact(signature, digestHash)
	if err != nil {
		return false
	}

	return pubKey.IsEqual(recoveredKey)
}

func CompressPubkey(pubKey *secp256k1.PublicKey) []byte {
	return pubKey.SerializeCompressed()
}

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [32]byte

func (h Hash) Bytes() []byte { return h[:] }
