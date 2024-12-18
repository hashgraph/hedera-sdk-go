package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
)

// _Ed25519PublicKey is an ed25519 public key.
type _Ed25519PublicKey struct {
	keyData []byte
}

// PublicKeyFromString recovers an _Ed25519PublicKey from its text-encoded representation.
func _Ed25519PublicKeyFromString(s string) (*_Ed25519PublicKey, error) {
	byt, err := hex.DecodeString(s)
	if err != nil {
		return &_Ed25519PublicKey{}, err
	}

	return _Ed25519PublicKeyFromBytes(byt)
}

// _Ed25519PublicKeyFromBytes constructs a known _Ed25519PublicKey from its text-encoded representation.
func _Ed25519PublicKeyFromBytes(bytes []byte) (*_Ed25519PublicKey, error) {
	length := len(bytes)
	switch length {
	case 32:
		return _Ed25519PublicKeyFromBytesRaw(bytes)
	case 44:
		return _Ed25519PublicKeyFromBytesDer(bytes)
	default:
		return &_Ed25519PublicKey{}, _NewErrBadKeyf("invalid public key length: %v bytes", len(bytes))
	}
}

// _Ed25519PublicKeyFromBytes constructs a known _Ed25519PublicKey from its text-encoded representation.
func _Ed25519PublicKeyFromBytesRaw(bytes []byte) (*_Ed25519PublicKey, error) {
	if bytes == nil {
		return &_Ed25519PublicKey{}, errByteArrayNull
	}

	if len(bytes) != ed25519.PublicKeySize {
		return &_Ed25519PublicKey{}, _NewErrBadKeyf("invalid public key length: %v bytes", len(bytes))
	}

	return &_Ed25519PublicKey{
		keyData: bytes,
	}, nil
}

func _Ed25519PublicKeyFromBytesDer(bytes []byte) (*_Ed25519PublicKey, error) {
	ed25519OID := asn1.ObjectIdentifier{1, 3, 101, 112}

	publicKeyInfo := struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}{}

	_, err := asn1.Unmarshal(bytes, &publicKeyInfo)
	if err != nil {
		return nil, err
	}

	if !publicKeyInfo.Algorithm.Algorithm.Equal(ed25519OID) {
		return nil, errors.New("invalid algorithm identifier, expected Ed25519")
	}

	if len(publicKeyInfo.PublicKey.Bytes) != 32 {
		return nil, errors.New("invalid public key length, expected 32 bytes")
	}

	var pk _Ed25519PublicKey
	pk.keyData = publicKeyInfo.PublicKey.Bytes

	return &pk, nil
}

func (pk _Ed25519PublicKey) _StringDer() string {
	return hex.EncodeToString(pk._BytesDer())
}

// String returns the text-encoded representation of the _Ed25519PublicKey.
func (pk _Ed25519PublicKey) _StringRaw() string {
	return hex.EncodeToString(pk.keyData)
}

// _Bytes returns the byte slice representation of the _Ed25519PublicKey.
func (pk _Ed25519PublicKey) _Bytes() []byte {
	return pk.keyData
}

// _Bytes returns the byte slice representation of the _Ed25519PublicKey.
func (pk _Ed25519PublicKey) _BytesRaw() []byte {
	return pk.keyData
}

func (pk _Ed25519PublicKey) _BytesDer() []byte {
	type PublicKeyInfo struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}

	ed25519OID := asn1.ObjectIdentifier{1, 3, 101, 112}
	publicKeyInfo := PublicKeyInfo{
		Algorithm: pkix.AlgorithmIdentifier{
			Algorithm: ed25519OID,
		},
		PublicKey: asn1.BitString{
			Bytes:     pk.keyData,
			BitLength: len(pk.keyData) * 8,
		},
	}

	derBytes, err := asn1.Marshal(publicKeyInfo)
	if err != nil {
		return nil
	}

	return derBytes
}

func (pk _Ed25519PublicKey) _ToProtoKey() *services.Key {
	return &services.Key{Key: &services.Key_Ed25519{Ed25519: pk.keyData}}
}

func (pk _Ed25519PublicKey) String() string {
	return pk._StringRaw()
}

func (pk _Ed25519PublicKey) _ToSignaturePairProtobuf(signature []byte) *services.SignaturePair {
	return &services.SignaturePair{
		PubKeyPrefix: pk.keyData,
		Signature: &services.SignaturePair_Ed25519{
			Ed25519: signature,
		},
	}
}

func (pk _Ed25519PublicKey) _Verify(message []byte, signature []byte) bool {
	return ed25519.Verify(pk._Bytes(), message, signature)
}

func (pk _Ed25519PublicKey) _VerifyTransaction(tx *Transaction[TransactionInterface]) bool {
	if tx.signedTransactions._Length() == 0 {
		return false
	}

	_, _ = tx._BuildAllTransactions()

	for _, value := range tx.signedTransactions.slice {
		tx := value.(*services.SignedTransaction)
		found := false
		for _, sigPair := range tx.SigMap.GetSigPair() {
			if bytes.Equal(sigPair.GetPubKeyPrefix(), pk._Bytes()) {
				found = true
				if !pk._Verify(tx.BodyBytes, sigPair.GetEd25519()) {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}
