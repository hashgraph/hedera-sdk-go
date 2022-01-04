package hedera

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
)

// _ECDSAPrivateKey is an Key_ECDSASecp256K1 private key.
type _ECDSAPrivateKey struct {
	*ecdsa.PrivateKey
}

func _GenerateECDSAPrivateKey() (*_ECDSAPrivateKey, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	return &_ECDSAPrivateKey{
		key,
	}, nil
}

func _ECDSAPrivateKeyFromBytes(byt []byte) (*_ECDSAPrivateKey, error) {
	length := len(byt)
	switch length {
	case 32:
		return _ECDSAPrivateKeyFromBytesRaw(byt)
	case 50:
		return _ECDSAPrivateKeyFromBytesDer(byt)
	default:
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}
}

func _ECDSAPrivateKeyFromBytesRaw(byt []byte) (*_ECDSAPrivateKey, error) {
	length := len(byt)
	if length != 32 {
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}

	b := bytes.Trim(byt, "\x00")

	d := new(big.Int)
	d.SetBytes(b)

	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = d

	x, y := elliptic.P256().ScalarBaseMult(b)
	privateKey.Curve = elliptic.P256()
	privateKey.X = x
	privateKey.Y = y

	return &_ECDSAPrivateKey{
		privateKey,
	}, nil
}

func _ECDSAPrivateKeyFromBytesDer(byt []byte) (*_ECDSAPrivateKey, error) {
	given := hex.EncodeToString(byt)

	result := strings.ReplaceAll(given, _ECDSAPrivateKeyPrefix, "")
	decoded, err := hex.DecodeString(result)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	if len(decoded) != 32 {
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}

	b := bytes.Trim(decoded, "\x00")

	d := new(big.Int)
	d.SetBytes(b)

	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = d

	x, y := elliptic.P256().ScalarBaseMult(b)
	privateKey.Curve = elliptic.P256()
	privateKey.X = x
	privateKey.Y = y

	return &_ECDSAPrivateKey{
		privateKey,
	}, nil
}

func _ECDSAPrivateKeyFromString(s string) (*_ECDSAPrivateKey, error) {
	b, err := hex.DecodeString(strings.ToLower(s))
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	return _ECDSAPrivateKeyFromBytes(b)
}

func (sk *_ECDSAPrivateKey) _PublicKey() *_ECDSAPublicKey {
	if sk.PrivateKey.Y == nil && sk.PrivateKey.X == nil {
		b := sk.D.Bytes()
		x, y := crypto.S256().ScalarBaseMult(b)
		sk.X = x
		sk.Y = y
		return &_ECDSAPublicKey{
			&ecdsa.PublicKey{
				Curve: crypto.S256(),
				X:     x,
				Y:     y,
			},
		}
	}

	return &_ECDSAPublicKey{
		&ecdsa.PublicKey{
			Curve: sk.Curve,
			X:     sk.X,
			Y:     sk.Y,
		},
	}
}

func (sk _ECDSAPrivateKey) _Sign(message []byte) []byte {
	hash := crypto.Keccak256Hash(message)
	sig, err := crypto.Sign(hash.Bytes(), sk.PrivateKey)
	if err != nil {
		panic(err)
	}

	// signature returned has a ecdsa recovery byte at the end,
	// need to remove it for verification to work.
	return sig[:len(sig)-1]
}

func (sk _ECDSAPrivateKey) _BytesRaw() []byte {
	privateKey := make([]byte, 32)
	temp := sk.D.Bytes()
	copy(privateKey[32-len(temp):], temp)

	return privateKey
}

func (sk _ECDSAPrivateKey) _BytesDer() []byte {
	prefix, _ := hex.DecodeString(_ECDSAPrivateKeyPrefix)
	return append(prefix, sk._BytesRaw()...)
}

func (sk _ECDSAPrivateKey) _StringDer() string {
	return fmt.Sprint(hex.EncodeToString(sk._BytesDer()))
}

func (sk _ECDSAPrivateKey) _StringRaw() string {
	return fmt.Sprint(hex.EncodeToString(sk._BytesRaw()))
}

func (sk _ECDSAPrivateKey) _ToProtoKey() *services.Key {
	return sk._PublicKey()._ToProtoKey()
}

func (sk _ECDSAPrivateKey) _SignTransaction(transaction *Transaction) ([]byte, error) {
	transaction._RequireOneNodeAccountID()

	if transaction.signedTransactions._Length() == 0 {
		return make([]byte, 0), errTransactionRequiresSingleNodeAccountID
	}

	signature := sk._Sign(transaction.signedTransactions._GetSignedTransactions()[0].GetBodyBytes())

	publicKey := sk._PublicKey()
	if publicKey == nil {
		return []byte{}, errors.New("public key is nil")
	}

	wrappedPublicKey := PublicKey{
		ecdsaPublicKey: publicKey,
	}

	if transaction._KeyAlreadySigned(wrappedPublicKey) {
		return []byte{}, nil
	}

	transaction.transactions = _NewLockedSlice()
	transaction.publicKeys = append(transaction.publicKeys, wrappedPublicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		temp := transaction.signedTransactions._GetSignedTransactions()[index]

		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		_, err := transaction.signedTransactions._Set(index, temp)
		if err != nil {
			transaction.lockError = err
		}
	}

	return signature, nil
}
